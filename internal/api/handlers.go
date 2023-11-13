package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"vio/internal/processes"

	"github.com/gorilla/mux"
)

// swagger:operation  GET /api/geolocation/{ip_address} GetGeoLocation
// Get Geo Location by IP address.
// ---
// produces:
// - application/json
// parameters:
//   - in: path
//     name: ip_address
//     description: IP address
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: OK
//	  schema:
//	    $ref: '#/definitions/Location'
//	'404':
//	  description: Error
//	'500':
//	  description: Internal Server error
func GetGeoLocation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ipAddress := vars["ip_address"]
		if vars == nil {
			// try to extract IP address from request (run from tests).
			ipAddress = processes.ExtractIPAddress(r.URL.String())
		}

		if !processes.IsValidIPAddress(ipAddress) {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}
		location, err := processes.GetLocation(db, ipAddress)
		if err != nil {
			http.Error(w, "Failed to retrieve location", http.StatusInternalServerError)
			return
		}
		if len(location.IPAddress) == 0 {
			http.Error(w, "Location not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(location)
	}
}
