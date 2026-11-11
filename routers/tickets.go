package routers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type ticketRouter struct {
	db *sql.DB
}

type Ticket struct {
	ID       int       `json:"id"`
	Cliente  string    `json:"cliente"`
	Asunto   string    `json:"asunto"`
	Status   string    `json:"status"`
	CreadoEn time.Time `json:"creado_en"`
}

func NewTicketsRouter(db *sql.DB) *ticketRouter {
	return &ticketRouter{db: db}
}

func (tr *ticketRouter) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", tr.getTicketsHandler)
	mux.HandleFunc("GET /{id}", tr.getTicketByIdHandler)
	mux.HandleFunc("POST /", tr.saveTicketHandler)
	mux.HandleFunc("PUT /{id}", tr.updateTicketHandler)
	mux.HandleFunc("DELETE /{id}", tr.deleteTicketHandler)
	return mux
}

func (tr *ticketRouter) getTicketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := tr.db.QueryContext(ctx, "SELECT * FROM tickets")
	if err != nil {

		http.Error(w, "Error al obtener los tickets", http.StatusNotFound)
		return
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		http.Error(w, "Error al leer los tickets", http.StatusNotFound)
		return
	}
	var tickets []Ticket
	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(&ticket.ID, &ticket.Cliente, &ticket.Asunto, &ticket.Status, &ticket.CreadoEn); err != nil {
			http.Error(w, "Error al escanear los tickets", http.StatusNotFound)
			return
		}
		tickets = append(tickets, ticket)
	}
	if len(tickets) == 0 {
		http.Error(w, "No se encontraron tickets", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tickets)
}

func (tr *ticketRouter) getTicketByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	var ticket Ticket
	if err := tr.db.QueryRowContext(ctx, "SELECT * FROM tickets WHERE id = $1", id).Scan(&ticket.ID, &ticket.Cliente, &ticket.Asunto, &ticket.Status, &ticket.CreadoEn); err != nil {
		http.Error(w, "Error al obtener el ticket", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ticket)
}

func (tr *ticketRouter) saveTicketHandler(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		http.Error(w, "Error al decodificar el ticket", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := tr.db.QueryRowContext(ctx, "INSERT INTO tickets (cliente, asunto) VALUES ($1, $2) RETURNING id, creado_en", ticket.Cliente, ticket.Asunto).Scan(&ticket.ID, &ticket.CreadoEn)
	if err != nil {
		http.Error(w, "Error al guardar el ticket", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

func (tr *ticketRouter) updateTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var ticket Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		http.Error(w, "Error al decodificar el ticket", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := tr.db.ExecContext(ctx, "UPDATE tickets SET estado = $1WHERE id = $3", ticket.Status, id)
	if err != nil {
		http.Error(w, "Error al actualizar el ticket", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ticket)
}

func (tr *ticketRouter) deleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := tr.db.ExecContext(ctx, "DELETE FROM tickets WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error al eliminar el ticket", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
