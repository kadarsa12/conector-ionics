package main

import (
	"database/sql"
	"errors"
	"strconv"

	_ "github.com/lib/pq"
	go_ora "github.com/sijms/go-ora/v2"
)

type database struct {
	*sql.DB
}

// const (
// 	batchSize      = 500
// 	maxRetries     = 3
// 	retryDelayBase = time.Second * 2
// )

func oracle_connection(host *string, port *int, name, user, pass, sid *string) (*database, error) {
	if host == nil || port == nil || name == nil || user == nil || pass == nil {
		return nil, errors.New("missing required parameters")
	}

	logger.Info("Connecting to Oracle DB...")

	var urlOptions map[string]string

	if sid != nil {
		urlOptions = map[string]string{
			"SID": *sid,
		}
	}

	connStr := go_ora.BuildUrl(*host, *port, *name, *user, *pass, urlOptions)
	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &database{conn}, nil
}

func pgsql_connection(host *string, port *int, name, user, pass, sid *string) (*database, error) {
	if host == nil || port == nil || name == nil || user == nil || pass == nil {
		return nil, errors.New("missing required parameters")
	}

	logger.Info("Connecting to PostgreSQL DB...")

	// connStr := go_ora.BuildUrl(*host, *port, *name, *user, *pass, urlOptions)
	connStr := "host=" + *host + " port=" + strconv.Itoa(*port) + " user=" + *user + " password=" + *pass + " dbname=" + *name + " sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &database{conn}, nil
}

func sqlsrv_connection(host *string, port *int, name, user, pass, sid *string) (*database, error) {
	if host == nil || port == nil || name == nil || user == nil || pass == nil {
		return nil, errors.New("missing required parameters")
	}

	logger.Info("Connecting to SQLServer DB...")

	var urlOptions map[string]string

	if sid != nil {
		urlOptions = map[string]string{
			"SID": *sid,
		}
	}

	connStr := go_ora.BuildUrl(*host, *port, *name, *user, *pass, urlOptions)
	conn, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, err
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &database{conn}, nil
}

func (d *database) get_data(batchSize, offset, customerID int, initialDate string) ([]map[string]interface{}, error) {
	query :=
		`select distinct
	$3 as api_customer_id,
	mpp.s_nome as nomeplanta,
	mt.i_cod_m_pessoa_empresa as i_cod_m_pessoa_empresa,
	mt.i_cod_m_transacao as i_cod_m_transacao,
	mt.dt_transa as dt_transa,
	ma.s_hora_inicial as s_hora_transa,
	ma.s_hora_final as s_hora_trans_final,
	mt.i_cod_b_tipo_transa as i_cod_b_tipo_transa,
	rvkh.s_result_valid_km_hr as s_result_valid_km,
	rvkh1.s_result_valid_km_hr as s_result_valid_hr,
	mdkh.s_motivo_digit_km_hr,
	m_item.s_desc as s_desc_item,
	m_item.i_cod as i_cod_item,
	pe_fornec.s_cod_pessoa as s_cod_fornecedor,
	pe_fornec.s_nome as s_nome_fornecedor,
	(case
		when bcss.i_cod_bs_tipo_un_vol_apres = 2
   then s_tipo_un_volume
		else b_unidade.s_desc
	end) as s_desc_unidade_medida,
	mt.i_cod_m_turno_log as i_cod_m_turno_log,
	(case
		when (bdfc.i_cod_bs_tipo_un_dist_persist = 1
		and bdfc.i_cod_bs_tipo_un_dist_apresent = 1)
		or (bdfc.i_cod_bs_tipo_un_dist_persist = 2
		and bdfc.i_cod_bs_tipo_un_dist_apresent = 2)
   then ma.f_km_atual
		when (bdfc.i_cod_bs_tipo_un_dist_persist = 1
		and bdfc.i_cod_bs_tipo_un_dist_apresent = 2)
   then ma.f_km_atual / bdfc.f_fator_conversao
		when (bdfc.i_cod_bs_tipo_un_dist_persist = 2
		and bdfc.i_cod_bs_tipo_un_dist_apresent = 1)
   then ma.f_km_atual * bdfc.f_fator_conversao
	end) as f_km_odometro,
	ma.f_horimetro_atual as i_ult_horimetro_informado,
	dmt.i_cod_m_veiculo as i_cod_m_veiculo,
	dmt.i_cod_r_unidade_item as i_cod_r_unidade_item,
	dmt.i_cod_m_bico_combust as i_cod_m_bico_combust,
	(case
		when (bvfc.i_cod_bs_tipo_un_vol_persist = 1
		and bvfc.i_cod_bs_tipo_un_vol_apresent = 1)
		or (bvfc.i_cod_bs_tipo_un_vol_persist = 2
		and bvfc.i_cod_bs_tipo_un_vol_apresent = 2)
   then dmt.f_quantidade
		when (bvfc.i_cod_bs_tipo_un_vol_persist = 1
		and bvfc.i_cod_bs_tipo_un_vol_apresent = 2)
   then dmt.f_quantidade / bvfc.f_fator_conversao
		when (bvfc.i_cod_bs_tipo_un_vol_persist = 2
		and bvfc.i_cod_bs_tipo_un_vol_apresent = 1)
   then dmt.f_quantidade * bvfc.f_fator_conversao
	end) as f_quantidade,
	coalesce(ma.f_preco_integracao,
	ma.f_preco_automacao) as f_preco_automacao,
	ma.f_valor as f_valor,
	ma.f_encerrante_inicial as f_encerrante_inicial,
	ma.f_encerrante_final as f_encerrante_final,
	dmt.i_cod_m_abastecimento as i_cod_m_abastecimento,
	mbc.s_desc as s_desc_bico,
	mbc.s_cod_ident_teclado as s_cod_ident_teclado_bico,
	b_bomba.s_desc as s_desc_bomba,
	mt.i_cod_bba_usuario as i_cod_bba_usuario,
	mt.i_cod_bba_usuario_motorista as i_cod_bba_usuario_motorista,
	bu_usuario.s_usuario_nome as s_desc_usuario,
	bu_motorista.s_usuario_nome as s_desc_motorista,
	mv.s_desc as s_desc_veiculo,
	mv.s_desc_red as s_desc_red_veiculo,
	mv.s_placa as s_placa,
	mv.s_cod_ident_teclado as s_cod_teclado_veiculo,
	mv.s_cod_sigla,
	mv.s_ordem_estatistica,
	mv.s_elemento_pep,
	mv.i_cod_bs_tipo_telemetria,
	coalesce(bcetmt.s_plano_cc_etapas,
	bcetve.s_plano_cc_etapas) as desc_centro_de_custo,
	coalesce(bcetmt.s_cod_b_plano_cc_etapas,
	bcetve.s_cod_b_plano_cc_etapas) as scod_centro_de_custo,
	b_sigla_veiculo.s_sigla as s_sigla,
	mt.i_cod_b_tipo_transa as i_cod_tipo_transa,
	bpv.s_desc as s_desc_plano_veic,
	mv.s_observacao as s_observacao,
	ve_pe.s_nome as s_nome_empresa_veiculo,
	pe_terceiro.s_cod_pessoa as s_cod_terceiro,
	pe_terceiro.s_nome as s_nome_terceiro,
	pe_terceiro_digitado.s_cod_pessoa as s_cod_terceiro_digitado,
	pe_terceiro_digitado.s_nome as s_nome_terceiro_digitado,
	mt.s_cod_veic_digitado,
	bi.s_cod_pista,
	plo.s_plano_operacao,
	plo.s_cod_plano_operacao,
	pec.s_nome as S_DESC_COTISTA,
	pec.i_cod_m_pessoa as codcotista,
	b_logradouro.s_nome as s_nome_logradouro,
	mt.dt_estorno,
	mt.i_cod_transa_origem_estorno,
	tis.s_tipo_integ_saaf,
	met.s_cod_ret_erp,
	met.s_desc_ret_erp,
	(case
		when met.i_cod_m_transacao is not null then
 'EXPORTADO'
		else 'NÃO EXPORTADO'
	end) as sExportacao,
	met.i_cod_m_transacao as transacao_exportada,
	tpt.s_desc as desc_bs_tipo_transa,
	mt.i_cod_b_safra,
	bsf.s_desc_safras as desc_safras,
	mt.i_cod_b_tipo_cultura,
	btc.s_desc_tipo_cultura as desc_tipo_cultura,
	mt.i_cod_b_tipo_atividade,
	bta.s_desc_tipo_atividade as desc_tipo_atividade,
	bla.s_id as s_id_local_armaz,
	btav.s_tipo_autoriz_veic,
	dmt.i_cod_d_m_transacao,
	t.s_desc as desc_Turno,
	btic.s_tipo_interface_comm,
	mmc.f_saldo_cota_corrente,
	bla.s_desc as descTanqueAbast,
	mcb.s_mid,
	blacomb.s_desc as descComboio,
	cab.i_esn_terminal,
	blamestre.s_desc as descComboioMestre,
	cab.dt_data_aquisicao,
	cab.s_hora_aquisicao,
	tau.s_tipo_autoriz_user,
	eacb.s_evento_abast,
	case
		when eacb.i_cod_bs_event_abast_cb = 13
     then 'SIM'
		else 'NÃO'
	end as senhadiaria,
	mt.dt_data_comp_bordo,
	hp_local.i_cod_h_pessoa_planta as iCodPessoaLocal,
	mod.s_desc as s_desc_modelo
from
	m_transacao mt
join d_m_transacao dmt
  on
	dmt.i_cod_m_transacao = mt.i_cod_m_transacao
	and dmt.i_cod_m_pessoa = mt.i_cod_m_pessoa_empresa
join m_abastecimento ma on
	ma.i_cod_m_abastecimento = dmt.i_cod_m_abastecimento
	and ma.i_cod_m_pessoa = dmt.i_cod_m_pessoa
	and ma.i_cod_bs_tipo_abastec in (1, 3, 8, 7)
	and (ma.i_cod_bs_event_abast_cb in (0, 1, 2, 3, 4, 6, 7, 10, 12, 13, 14, 15, 20, 21, 22, 24, 23)
		or ma.i_cod_bs_event_abast_cb is null)
left join bs_event_abast_cb eacb
    on
	eacb.i_cod_bs_event_abast_cb = ma.i_cod_bs_event_abast_cb
left join m_bico_combust mbc on
	mbc.i_cod = dmt.i_cod_m_bico_combust
	and mbc.i_cod_m_pessoa = dmt.i_cod_m_pessoa
left join m_tanque_logico on
	m_tanque_logico.i_cod = mbc.i_cod_m_tanque_logico
	and m_tanque_logico.i_cod_m_pessoa = mbc.i_cod_m_pessoa
left join b_local_armaz bla on
	(m_tanque_logico.i_cod_b_local_armaz = bla.i_cod
		and bla.i_cod_m_pessoa = m_tanque_logico.i_cod_m_pessoa)
left join bs_tipo_autoriz_veic btav
  on
	btav.i_cod_bs_tipo_autoriz_veic = ma.i_cod_bs_tipo_autoriz_veic
left join b_logradouro on
	b_logradouro.i_cod = ma.i_cod_b_logradouro
left join m_veiculo mv
     on
	mv.i_cod = dmt.i_cod_m_veiculo
left join b_modelo mod on
	(mv.i_cod_b_mod = mod.i_cod)
left join b_sigla_veiculo on
	(b_sigla_veiculo.i_cod_b_sigla = mv.i_cod_b_sigla_veiculo)
left join b_plano_cc_etapas bcetmt on
	(mt.i_cod_b_plano_cc_etapas = bcetmt.i_cod_b_plano_cc_etapas)
left join b_plano_cc_etapas bcetve on
	(mv.i_cod_b_plano_cc_etapas = bcetve.i_cod_b_plano_cc_etapas)
left join r_veiculo_pessoa ve_vp
     on
	ve_vp.i_cod_m_veiculo = mv.i_cod
left join m_pessoa ve_pe
     on
	ve_pe.i_cod_m_pessoa = coalesce(mt.i_cod_pessoa_r_veiculo,
	ve_vp.i_cod_m_pessoa)
left join b_plano_veiculo bpv on
	bpv.i_cod = mv.i_cod_b_plano_veiculo
left join bba_usuario bu_usuario on
	bu_usuario.i_cod = mt.i_cod_bba_usuario
left join bba_usuario bu_motorista on
	bu_motorista.i_cod = mt.i_cod_bba_usuario_motorista
left join m_pessoa pe_terceiro on
	pe_terceiro.i_cod_m_pessoa = dmt.i_cod_m_pessoa_terceiro
left join m_pessoa pe_terceiro_digitado on
	pe_terceiro_digitado.i_cod_m_pessoa = mt.i_cod_m_pessoa_terceiro
left join b_intercomm bi on
	(bi.i_cod_b_intercomm = mt.i_cod_b_intercomm
		and bi.i_cod_m_pessoa = mt.i_cod_m_pessoa_empresa)
left join m_mov_cota mmc on
		mmc.i_cod_m_transacao = mt.i_cod_m_transacao
	and mmc.i_cod_m_pessoa = mt.i_cod_m_pessoa_empresa
left join b_plano_operacao plo on
	(plo.i_cod_b_plano_operacao = dmt.i_cod_b_plano_operacao)
inner join b_config_sistema_saaf bcss on
	(1 = 1)
inner join bs_tipo_un_volume bv on
	(bcss.i_cod_bs_tipo_un_vol_apres = bv.i_cod_bs_tipo_un_volume)
inner join bs_tipo_un_distancia bd on
	(bcss.i_cod_bs_tipo_un_dist_apres = bd.i_cod_bs_tipo_un_distancia)
left join bs_tipo_transa tpt on
	tpt.i_cod = mt.i_cod_b_tipo_transa
left join b_safras bsf on
	bsf.i_cod_b_safras = mt.i_cod_b_safra
left join b_tipo_cultura btc on
	btc.i_cod_b_tipo_cultura = mt.i_cod_b_tipo_cultura
left join b_tipo_atividade bta on
	bta.i_cod_b_tipo_atividade = mt.i_cod_b_tipo_atividade
left join m_pessoa pec on                               
		pec.i_cod_m_pessoa = mmc.i_cod_m_pessoa_cotista
left join h_pessoa_planta hp_ext on
	(hp_ext.i_cod_h_pessoa_planta = mt.i_cod_m_pessoa_empresa)
left join h_pessoa_planta hp_local on
	(hp_local.ib_local = 1)
left join bs_tipo_un_volume_fc bvfc on
	(bcss.i_cod_bs_tipo_un_vol_apres = bvfc.i_cod_bs_tipo_un_vol_apresent
		and bvfc.i_cod_bs_tipo_un_vol_persist = hp_ext.i_cod_bs_tipo_un_vol_persist )
left join bs_tipo_un_distancia_fc bdfc on
	(bcss.i_cod_bs_tipo_un_dist_apres = bdfc.i_cod_bs_tipo_un_dist_apresent
		and bdfc.i_cod_bs_tipo_un_dist_persist = hp_ext.i_cod_bs_tipo_un_dist_persist )
left join bs_result_valid_km_hr rvkh on
	(mt.i_cod_bs_result_valid_km = rvkh.i_cod_bs_result_valid_km_hr)
left join bs_result_valid_km_hr rvkh1 on
	(mt.i_cod_bs_result_valid_hr = rvkh1.i_cod_bs_result_valid_km_hr)
left join bs_motivo_digit_km_hr mdkh on
	(mt.i_cod_bs_motivo_digit_km_hr = mdkh.i_cod_bs_motivo_digit_km_hr)
left join m_exp_transacao met 
  on
	met.i_cod_m_transacao = mt.i_cod_m_transacao
	and met.i_cod_m_pessoa_transa = mt.i_cod_m_pessoa_empresa
left join bs_tipo_integ_saaf tis
  on
	met.i_cod_bs_tipo_integ_saaf = tis.i_cod_bs_tipo_integ_saaf
left join b_bomba on
	b_bomba.i_cod = mbc.i_cod_b_bomba
	and b_bomba.i_cod_m_pessoa = mbc.i_cod_m_pessoa
left join m_pessoa mpp on
	mpp.i_cod_m_pessoa = mt.i_cod_m_pessoa_empresa
left join r_unidade_item on
	r_unidade_item.i_cod = dmt.i_cod_r_unidade_item
left join b_unidade on
	b_unidade.i_cod = r_unidade_item.i_cod_b_unidade
left join m_item on
	m_item.i_cod = r_unidade_item.i_cod_m_item
left join (
	select
		MAX(i_cod_h_pessoa_juridica) as i_cod_h_pessoa_juridica,
		i_cod_m_item
	from
		r_for_item
	group by
		i_cod_m_item) rfi on
	m_item.i_cod = rfi.i_cod_m_item
left join m_pessoa pe_fornec on
	pe_fornec.i_cod_m_pessoa = rfi.i_cod_h_pessoa_juridica
left join m_turno_log tlog
  on
	tlog.i_cod = mt.i_cod_m_turno_log
	and tlog.i_cod_m_pessoa = mt.i_cod_m_pessoa_empresa
left join b_turno t
  on
	t.i_cod = tlog.i_cod_b_turno
left join bs_tipo_interface_comm btic
  on
	btic.i_cod_bs_tipo_interface_comm = mt.i_cod_bs_tipo_interface_comm
left join bs_tipo_autoriz_user tau
  on
	tau.i_cod_bs_tipo_autoriz_user = mt.i_cod_bs_tipo_autoriz_user
left join m_cab_database_cb cab
  on
	cab.i_cod_m_cab_database_cb = mt.i_cod_m_cab_database_cb
	and cab.i_cod_m_pessoa = mt.i_cod_m_pessoa_empresa
left join m_comboio mcbmestre
  on
	mcbmestre.i_cod_m_comboio = cab.i_cod_m_combio
	and mcbmestre.i_cod_m_pessoa = cab.i_cod_m_pessoa
left join m_tanque_logico mtlmestre
  on
	mtlmestre.i_cod = mcbmestre.i_cod_m_tanque_logico
	and mtlmestre.i_cod_m_pessoa = mcbmestre.i_cod_m_pessoa
left join b_local_armaz blamestre
  on
	blamestre.i_cod = mtlmestre.i_cod_b_local_armaz
	and blamestre.i_cod_m_pessoa = mtlmestre.i_cod_m_pessoa
left join m_comboio mcb
  on
	mcb.i_cod_m_comboio = ma.i_cod_m_comboio
	and mcb.i_cod_m_pessoa = ma.i_cod_m_pessoa
left join m_tanque_logico mtlcomb
  on
	mtlcomb.i_cod = mcb.i_cod_m_tanque_logico
	and mtlcomb.i_cod_m_pessoa = mcb.i_cod_m_pessoa
left join b_local_armaz blacomb
  on
	blacomb.i_cod = mtlcomb.i_cod_b_local_armaz
	and blacomb.i_cod_m_pessoa = mtlcomb.i_cod_m_pessoa
where
	mt.i_cod_b_tipo_transa in (101, 112, 117)
	and (1 = 1)
	and (bpv.i_cod is null
		or (1 = 1) )
	and (mv.i_cod_b_plano_veiculo is null
		or (1 = 1) )
	and mt.dt_estorno is null
	and mt.i_cod_transa_origem_estorno is null
	and mt.dt_transa > cast($4 as date)
	and dmt.f_quantidade > 0
order by i_cod_m_transacao desc 
LIMIT $1 OFFSET $2`

	rows, err := d.Query(query, batchSize, offset, customerID, initialDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Obter os nomes das colunas
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var records []map[string]interface{}

	for rows.Next() {
		// Cria um slice para armazenar os valores das colunas
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Faz o scan dos valores para o slice
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Cria um mapa para o registro atual
		record := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]

			// Trata valores do tipo []byte (tipicamente para colunas VARCHAR, TEXT, etc.)
			if b, ok := val.([]byte); ok {
				record[colName] = string(b)
			} else {
				record[colName] = val
			}
		}

		records = append(records, record)
	}

	return records, nil
}

// func (d *database) get_data2(batchSize, offset) ([]string, error) {
// 	offset := 0
// 	// rateLimiter := ratelimit.New(5) // Limita a 5 requisições por segundo

// 	for {
// 		records, err := fetchRecords(d, offset, batchSize)
// 		if err != nil {
// 			log.Printf("Erro ao buscar registros: %v", err)
// 			break
// 		}

// 		if len(records) == 0 {
// 			log.Println("Nenhum registro encontrado, processo concluído.")
// 			break
// 		}

// 		// Envia o lote para a API usando goroutine
// 		// rateLimiter.Take()
// 		// go func(records []map[string]interface{}) {
// 		// 	if err := sendToAPIWithRetry(records); err != nil {
// 		// 		// log.Printf("Erro ao enviar registros para a API: %v", err)
// 		// 		logger.Info("Erro ao enviar registros para a API: %v", err)
// 		// 	}
// 		// }(records)

// 		// offset += batchSize
// 		// time.Sleep(time.Millisecond * 500)
// 	}

// 	// Aguarda as goroutines finalizarem
// 	// Isso é um simplificação; em produção, considere usar sync.WaitGroup
// 	// time.Sleep(time.Minute * 5)
// 	// return nil, nil
// }

// func sendToAPIWithRetry(records []map[string]interface{}) error {
// 	var err error
// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		err = sendToAPI(records)
// 		if err == nil {
// 			// log.Printf("Lote de %d registros enviado com sucesso.", len(records))
// 			logger.Info("Lote de %d registros enviado com sucesso.", len(records))
// 			return nil
// 		}
// 		// log.Printf("Tentativa %d falhou: %v. Retentando em %v...", attempt, err, retryDelayBase*time.Duration(attempt))
// 		logger.Info("Tentativa %d falhou: %v. Retentando em %v...", attempt, err, retryDelayBase*time.Duration(attempt))

// 		time.Sleep(retryDelayBase * time.Duration(attempt))
// 	}
// 	return fmt.Errorf("falha ao enviar registros após %d tentativas: %v", maxRetries, err)
// }

// func sendToAPI(records []map[string]interface{}) error {
// 	client := &http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	data, err := json.Marshal(records)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest("POST", "http://127.0.0.1:3333/v1/writer", bytes.NewReader(data))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// Envia a requisição
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Verifica o status da resposta
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("status code inesperado: %d", resp.StatusCode)
// 	}

// 	return nil
// }
