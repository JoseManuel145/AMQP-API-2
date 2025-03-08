package infraestructure

import (
	"fmt"
	"log"
	"rabbitConsumer/src/core"
	_ "rabbitConsumer/src/report/domain"
	"rabbitConsumer/src/report/domain/entities"
)

type MySQL struct {
	conn *core.Conn_MySQL
}

func NewMySQL() *MySQL {
	conn := core.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (r *MySQL) Update(id int, title, content, status string) error {
	query := "UPDATE reports SET title = ?, content = ?, status = ? WHERE id = ?"
	_, err := r.conn.ExecutePreparedQuery(query, title, content, status, id)
	return err
}

func (r *MySQL) GetAll() ([]entities.Report, error) {
	query := "SELECT * FROM reports"
	rows := r.conn.FetchRows(query)
	defer rows.Close()

	var reports []entities.Report
	for rows.Next() {
		var report entities.Report
		if err := rows.Scan(&report.ID, &report.Title, &report.Content, &report.Status); err != nil {
			return nil, fmt.Errorf("error al escanear la fila: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando sobre las filas: %w", err)
	}
	return reports, nil
}
