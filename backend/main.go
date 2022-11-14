package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	// "io"
	// "log"
)

type Neo4jConfiguration struct {
	Url      string
	Username string
	Password string
	Database string
}

type GarphologyResponse struct {
	Nodes []Node `json:"nodes"`
	Links []Edge `json:"edges"`
}
type NodeRequest struct {
	Method string `json:"method"`
	Node   `json:"payload"`
}

type NodeResult struct {
	Success int `json:"success"`
}
type Node struct {
	Identity   int64                  `json:"key"`
	Properties map[string]interface{} `json:"attributes,omitempty"`
}

type Edge struct {
	Identity int64 `json:"key"`
	// Type       string                 `json:"type"`
	Start      int64                  `json:"source"`
	End        int64                  `json:"target"`
	Properties map[string]interface{} `json:"attributes,omitempty"`
}

func parseConfiguration() *Neo4jConfiguration {

	return &Neo4jConfiguration{
		Url:      "neo4j://localhost:7687",
		Username: "neo4j",
		Password: "Northwind",
		Database: "neo4j",
	}
}

func (nc *Neo4jConfiguration) newDriver() (neo4j.Driver, error) {
	return neo4j.NewDriver(nc.Url, neo4j.BasicAuth(nc.Username, nc.Password, ""))
}

func unsafeClose(closeable io.Closer) {
	if err := closeable.Close(); err != nil {
		log.Fatal(fmt.Errorf("could not close resource: %w", err))
	}
}

func parseLimit(req *http.Request) int {
	limits := req.URL.Query()["limit"]
	limit := 50
	if len(limits) > 0 {
		var err error
		if limit, err = strconv.Atoi(limits[0]); err != nil {
			limit = 50
		}
	}
	return limit
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	_, file, _, _ := runtime.Caller(0)
	page := filepath.Join(filepath.Dir(file), "public", "index.html")
	fmt.Printf("Serving HTML file %s\n", page)
	if body, err := ioutil.ReadFile(page); err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		_, _ = w.Write(body)
	}
}

func graphHandler(driver neo4j.Driver, database string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		session := driver.NewSession(neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeRead,
			DatabaseName: database,
		})
		defer unsafeClose(session)

		limit := 1000000
		query_nodes := `MATCH (n) RETURN labels(n) as l, ID(n) as id, properties(n) as p `
		query_edges := `MATCH (sr)-[r]->(er) RETURN ID(r) as rid,  properties(r) as rprops, type(r) as rtype, ID(sr) as srid, ID(er) as erid`
		fmt.Println(query_nodes)
		fmt.Println(query_edges)
		d3Resp, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			records_node, err := tx.Run(query_nodes, map[string]interface{}{"limit": limit})
			if err != nil {
				return nil, err
			}

			result := GarphologyResponse{}
			for records_node.Next() {
				record := records_node.Record()
				// fmt.Printf("%s\n", record)
				labels, _ := record.Get("l")
				identity, _ := record.Get("id")
				properties, _ := record.Get("p")
				node := Node{}
				if rec, ok := identity.(int64); ok {
					node.Identity = rec
				} else {
					fmt.Printf("identity not a int: %v\n", identity)
				}

				if rec_p, ok := properties.(map[string]interface{}); ok {
					if rec_l, ok := labels.([]interface{}); ok {
						rec_p["labels"] = rec_l
					} else {
						fmt.Printf("labels not a []interface{}: %v\n", labels)
					}
					node.Properties = rec_p
					// for key, val := range rec {
					// 	fmt.Println(key, val)
					// }
				} else {
					fmt.Printf("record not a map[string]interface{}: %v\n", record)
				}

				result.Nodes = append(result.Nodes, node)

			}

			// get all eages
			records_edge, err := tx.Run(query_edges, map[string]interface{}{"limit": limit})
			if err != nil {
				return nil, err
			}
			for records_edge.Next() {
				record := records_edge.Record()
				// fmt.Printf("%s\n", record)

				identity, _ := record.Get("rid")
				properties, _ := record.Get("rprops")
				rtype, _ := record.Get("rtype")
				startID, _ := record.Get("srid")
				endID, _ := record.Get("erid")

				link := Edge{}
				if rec, ok := identity.(int64); ok {
					link.Identity = rec
				} else {
					fmt.Printf("identity not a int: %v\n", identity)
				}

				if rec, ok := startID.(int64); ok {
					link.Start = rec
				} else {
					fmt.Printf("startID not a int: %v\n", startID)
				}

				if rec, ok := endID.(int64); ok {
					link.End = rec
				} else {
					fmt.Printf("endID not a int: %v\n", endID)
				}

				if rec_p, ok := properties.(map[string]interface{}); ok {
					// add type of link to properties
					if rec_t, ok := rtype.(string); ok {
						rec_p["edge_type"] = rec_t
					} else {
						fmt.Printf("rec_t not a string: %v\n", rtype)
					}

					link.Properties = rec_p
					// for key, val := range rec {
					// 	fmt.Println(key, val)
					// }
				} else {
					fmt.Printf("record not a map[string]interface{}: %v\n", record)
				}

				result.Links = append(result.Links, link)

			}
			// fmt.Println(result)
			return result, nil
		})
		if err != nil {
			log.Println("error querying graph:", err)
			return
		}
		err = json.NewEncoder(w).Encode(d3Resp)
		if err != nil {
			log.Println("error writing graph response:", err)
		}
	}
}

func nodeHandler(driver neo4j.Driver, database string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		case "POST":
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			defer req.Body.Close()
			buf, err := io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			// fmt.Println(buf)
			var nReq NodeRequest
			json.NewDecoder(strings.NewReader(string(buf))).Decode(&nReq)
			session := driver.NewSession(neo4j.SessionConfig{
				AccessMode:   neo4j.AccessModeRead,
				DatabaseName: database,
			})
			defer unsafeClose(session)
			if nReq.Method == "add" {
				addNode(w, req, session, nReq.Node)
			} else if nReq.Method == "delete" {
				deleteNode(w, req, session, nReq.Node)
			} else if nReq.Method == "update" {
				updateNode(w, req, session, nReq.Node)
			} else {
				fmt.Println("to be continue ...")
			}
		}
	}
}

func addNode(w http.ResponseWriter, req *http.Request, session neo4j.Session, newnode Node) {
	NODE_LABEL := "labels"

	var buffer bytes.Buffer
	buffer.WriteString("CREATE (n")
	for key, val := range newnode.Properties {
		if key == NODE_LABEL {
			buffer.WriteString(":")
			if rec, ok := val.(string); ok {
				buffer.WriteString(rec)
			} else {
				log.Println("error: the type of 'label' is not string.")
			}
			break
		}
	}
	buffer.WriteString("{")
	index := 0
	for key, val := range newnode.Properties {
		if key == NODE_LABEL {
			continue
		}
		if index != 0 {
			buffer.WriteString(", ")
		}
		index++
		buffer.WriteString(key)
		buffer.WriteString(": ")
		if rec, ok := val.(string); ok {
			buffer.WriteString("'")
			buffer.WriteString(rec)
			buffer.WriteString("'")
		} else {
			log.Println("error: the type of attribute is not string.")
		}
	}
	buffer.WriteString("})")

	add_node_cypher := buffer.String()
	fmt.Println(add_node_cypher)
	addNodeResp, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(add_node_cypher, nil)
		if err != nil {
			return nil, err
		}
		// fmt.Println(result)
		var summary, _ = result.Consume()
		var addNodeResult NodeResult
		addNodeResult.Success = summary.Counters().NodesCreated()
		fmt.Println(addNodeResult)

		// return the number of nodes created.
		return addNodeResult, nil
	})
	if err != nil {
		log.Println("error adding node:", err)
		return
	}
	err = json.NewEncoder(w).Encode(addNodeResp)
	if err != nil {
		log.Println("error writing node response:", err)
	}
}

func deleteNode(w http.ResponseWriter, req *http.Request, session neo4j.Session, newnode Node) {
	nodeID := newnode.Identity
	addNodeCypher := `MATCH (n) WHERE ID(n) = $nodeID DELETE (n)`
	fmt.Println(addNodeCypher)
	addNodeResp, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(addNodeCypher, map[string]interface{}{
			"nodeID": nodeID,
		})
		if err != nil {
			return nil, err
		}
		fmt.Println(result)
		var summary, _ = result.Consume()
		var deleteNodeResult NodeResult
		deleteNodeResult.Success = summary.Counters().NodesDeleted()
		fmt.Println(deleteNodeResult)

		// return the number of nodes created.
		return deleteNodeResult, nil
	})
	if err != nil {
		log.Println("error adding node:", err)
		return
	}
	err = json.NewEncoder(w).Encode(addNodeResp)
	if err != nil {
		log.Println("error writing node response:", err)
	}
}

func updateNode(w http.ResponseWriter, req *http.Request, session neo4j.Session, newnode Node) {
	nodeID := newnode.Identity
	NODE_LABEL := "labels"

	var buffer bytes.Buffer
	buffer.WriteString("MATCH (n) WHERE ID(n) = $nodeID ")
	for key, val := range newnode.Properties {
		if key == NODE_LABEL {
			buffer.WriteString("SET n:")
			if rec, ok := val.(string); ok {
				buffer.WriteString(rec)
				buffer.WriteString(" ")
			} else {
				log.Println("error: the type of 'label' is not string.")
			}
		} else {
			buffer.WriteString("SET n.")
			buffer.WriteString(key)
			buffer.WriteString("=")
			if rec, ok := val.(string); ok {
				buffer.WriteString("'")
				buffer.WriteString(rec)
				buffer.WriteString("' ")
			} else {
				log.Println("error: the type of 'label' is not string.")
			}
		}

	}
	buffer.WriteString("RETURN (n)")
	updateNodeCypher := buffer.String()
	fmt.Println(updateNodeCypher)
	updateNodeResp, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(updateNodeCypher, map[string]interface{}{
			"nodeID": nodeID,
		})
		if err != nil {
			return nil, err
		}
		fmt.Println(result)
		var summary, _ = result.Consume()
		var updateNodeResult NodeResult

		if summary.Counters().ContainsUpdates() {
			updateNodeResult.Success = 1
		} else {
			updateNodeResult.Success = 0
		}

		fmt.Println(updateNodeResult)
		// return the number of nodes created.
		return updateNodeResult, nil
	})
	if err != nil {
		log.Println("error adding node:", err)
		return
	}
	err = json.NewEncoder(w).Encode(updateNodeResp)
	if err != nil {
		log.Println("error writing node response:", err)
	}
}
func main() {
	configuration := parseConfiguration()

	driver, err := configuration.newDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer unsafeClose(driver)
	serveMux := http.NewServeMux()
	// serveMux.HandleFunc("/", defaultHandler)
	serveMux.HandleFunc("/", graphHandler(driver, configuration.Database))
	serveMux.HandleFunc("/node", nodeHandler(driver, configuration.Database))
	fmt.Println(configuration)

	var port string
	var found bool
	if port, found = os.LookupEnv("PORT"); !found {
		port = "8083"
	}
	fmt.Printf("Running on port %s, database is at %s\n", port, configuration.Url)
	panic(http.ListenAndServe(":"+port, serveMux))
}
