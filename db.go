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

// func (d *database) get_invoices(keys []string) ([]invoice, error) {
// 	if len(keys) == 0 {
// 		return nil, nil
// 	}

// 	query := "SELECT * FROM FRETE_GESTAO.VSDI_DADOS_NFE_TMS WHERE NFE_KEY IN ("

// 	params := []any{}
// 	for i, key := range keys {
// 		if i > 0 {
// 			query += ","
// 		}
// 		query += ":" + strconv.Itoa(i+1)
// 		params = append(params, key)
// 	}
// 	query += ")"

// 	rows, err := d.Query(query, params...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var invoices []invoice
// 	for rows.Next() {
// 		var inv invoice
// 		err := rows.Scan(
// 			&inv.NfeKey,
// 			&inv.NfeCreatedAt,
// 			&inv.NfeNumber,
// 			&inv.NfeSequence,
// 			&inv.NfeNumPed,
// 			&inv.NfeTipo,
// 			&inv.NfeWeight,
// 			&inv.NfeCubage,
// 			&inv.NfeCubicWeight,
// 			&inv.NfeAmount,
// 			&inv.VolumeNfeAmount,
// 			&inv.VolumeNfeType,
// 			&inv.EmitterNfeName,
// 			&inv.EmitterNfeCnpj,
// 			&inv.EmitterAddressNfeZipcode,
// 			&inv.EmitterAddressNfeNumber,
// 			&inv.EmitterAddressNfeDistrict,
// 			&inv.EmitterAddressNfeCity,
// 			&inv.EmitterAddressNfeState,
// 			&inv.EmitterAddressNfeCountry,
// 			&inv.EmitterAddressNfeStreet,
// 			&inv.ReceiverNfeName,
// 			&inv.ReceiverNfeDocumentNumber,
// 			&inv.ReceiverAddressNfeStreet,
// 			&inv.ReceiverAddressNfeNumber,
// 			&inv.ReceiverAddressNfeZipcode,
// 			&inv.ReceiverAddressNfeDistrict,
// 			&inv.ReceiverAddressNfeCity,
// 			&inv.ReceiverAddressNfeState,
// 			&inv.ReceiverAddressNfeCountry,
// 			&inv.CarrierNfeName,
// 			&inv.CarrierNfeCnpj,
// 			&inv.CarrierNfeFreightMode,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		inv.InOrOut = "I"

// 		invoices = append(invoices, inv)
// 	}

// 	query = "SELECT * FROM FRETE_GESTAO.VSDI_DADOS_NFS_TMS WHERE NFE_KEY IN ("

// 	params = []any{}
// 	for i, key := range keys {
// 		if i > 0 {
// 			query += ","
// 		}
// 		query += ":" + strconv.Itoa(i+1)
// 		params = append(params, key)
// 	}
// 	query += ")"

// 	rows, err = d.Query(query, params...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var inv invoice
// 		err := rows.Scan(
// 			&inv.NfeKey,
// 			&inv.NfeCreatedAt,
// 			&inv.NfeNumber,
// 			&inv.NfeSequence,
// 			&inv.NfeNumPed,
// 			&inv.NfeTipo,
// 			&inv.NfeWeight,
// 			&inv.NfeCubage,
// 			&inv.NfeCubicWeight,
// 			&inv.NfeAmount,
// 			&inv.VolumeNfeAmount,
// 			&inv.VolumeNfeType,
// 			&inv.EmitterNfeName,
// 			&inv.EmitterNfeCnpj,
// 			&inv.EmitterAddressNfeZipcode,
// 			&inv.EmitterAddressNfeNumber,
// 			&inv.EmitterAddressNfeDistrict,
// 			&inv.EmitterAddressNfeCity,
// 			&inv.EmitterAddressNfeState,
// 			&inv.EmitterAddressNfeCountry,
// 			&inv.EmitterAddressNfeStreet,
// 			&inv.ReceiverNfeName,
// 			&inv.ReceiverNfeDocumentNumber,
// 			&inv.ReceiverAddressNfeStreet,
// 			&inv.ReceiverAddressNfeNumber,
// 			&inv.ReceiverAddressNfeZipcode,
// 			&inv.ReceiverAddressNfeDistrict,
// 			&inv.ReceiverAddressNfeCity,
// 			&inv.ReceiverAddressNfeState,
// 			&inv.ReceiverAddressNfeCountry,
// 			&inv.CarrierNfeName,
// 			&inv.CarrierNfeCnpj,
// 			&inv.CarrierNfeFreightMode,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		inv.InOrOut = "O"

// 		invoices = append(invoices, inv)
// 	}

// 	return invoices, nil
// }

// func (d *database) insert_status(ctes []Cte, status string, companyID int) error {
// 	if len(ctes) == 0 {
// 		return nil
// 	}

// 	now := time.Now()
// 	var err error

// 	for _, cte := range ctes {
// 		logger.Info("Inserting invoice status: '" + cte.Key + "' to " + status)
// 		_, err = d.Exec(`INSERT INTO FOCCO3I.TTMS_STATUS_CTE
// 			(ID, CREATED_BY, MODIFIED_BY, MODIFIED_ON, EMPR_ID, NRO_DOC, SERIE, CHAVE_ACESSO, STATUS) VALUES (FOCCO3I.SEQ_ID_TTMS_STATUS_CTE.NEXTVAL, 'FRETE_GESTAO', 'FRETE_GESTAO', :1, :2, :3, :4, :5, :6)`, now, companyID, cte.Number, cte.Serie, cte.Key, status)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (d *database) remove_status(keys []string, companyID int) error {
// 	if len(keys) == 0 {
// 		return nil
// 	}

// 	var err error
// 	for _, key := range keys {
// 		logger.Info("Removing invoice status: " + key)
// 		_, err = d.Exec(`DELETE FROM FOCCO3I.TTMS_STATUS_CTE WHERE CHAVE_ACESSO = '` + key + `' AND CREATED_BY = 'FRETE_GESTAO' AND EMPR_ID = ` + strconv.Itoa(companyID))
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// type invoice struct {
// 	NfeKey                     *string `json:"key"`
// 	NfeCreatedAt               *string `json:"created_at"`
// 	NfeNumber                  *int    `json:"number"`
// 	NfeSequence                *int    `json:"sequence"`
// 	NfeNumPed                  *int    `json:"num_ped"`
// 	NfeTipo                    *string `json:"tipo"`
// 	NfeWeight                  *string `json:"weight"`
// 	NfeCubage                  *string `json:"cubage"`
// 	NfeCubicWeight             *string `json:"cubic_weight"`
// 	NfeAmount                  *string `json:"amount"`
// 	VolumeNfeAmount            *int    `json:"volume_amount"`
// 	VolumeNfeType              *string `json:"volume_type"`
// 	EmitterNfeName             *string `json:"emitter_name"`
// 	EmitterNfeCnpj             *string `json:"emitter_doc"`
// 	EmitterAddressNfeZipcode   *int    `json:"emitter_zipcode"`
// 	EmitterAddressNfeNumber    *string `json:"emitter_number"`
// 	EmitterAddressNfeDistrict  *string `json:"emitter_district"`
// 	EmitterAddressNfeCity      *string `json:"emitter_city"`
// 	EmitterAddressNfeState     *string `json:"emitter_state"`
// 	EmitterAddressNfeCountry   *string `json:"emitter_country"`
// 	EmitterAddressNfeStreet    *string `json:"emitter_street"`
// 	ReceiverNfeName            *string `json:"receiver_name"`
// 	ReceiverNfeDocumentNumber  *string `json:"receiver_doc"`
// 	ReceiverAddressNfeStreet   *string `json:"receiver_street"`
// 	ReceiverAddressNfeNumber   *string `json:"receiver_number"`
// 	ReceiverAddressNfeZipcode  *string `json:"receiver_zipcode"`
// 	ReceiverAddressNfeDistrict *string `json:"receiver_district"`
// 	ReceiverAddressNfeCity     *string `json:"receiver_city"`
// 	ReceiverAddressNfeState    *string `json:"receiver_state"`
// 	ReceiverAddressNfeCountry  *string `json:"receiver_country"`
// 	CarrierNfeName             *string `json:"carrier_name"`
// 	CarrierNfeCnpj             *string `json:"carrier_doc"`
// 	CarrierNfeFreightMode      *string `json:"carrier_freightmode"`

// 	InOrOut string `json:"in_out"`
// }
