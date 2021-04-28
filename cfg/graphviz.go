package cfg

import (
	"github.com/awalterschulze/gographviz"
	"strconv"
)

func GraphViz(cfg *CFG, cfgFuncs map[string]*CFG) *gographviz.Graph {
	graphAst, _ := gographviz.ParseString(`digraph toplevel {}`)
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}

	buildGraph(graph, cfg, "toplevel")

	for id, cfg := range cfgFuncs {
		subgraph := "cluster_" + id
		err := graph.AddSubGraph("toplevel", subgraph, map[string]string{"label": id})
		if err != nil {
			panic(err)
		}
		buildGraph(graph, cfg, subgraph)
	}

	return graph
}

func buildGraph(graph *gographviz.Graph, cfg *CFG, parentGraph string) {
	cfg.Visit(func(node *Node) {
		currLabel := strconv.Itoa(node.Label)
		exprLabel := strconv.Quote(node.Expr.String())
		err := graph.AddNode(parentGraph, currLabel, map[string]string{"label": exprLabel})
		if err != nil {
			panic(err)
		}

		if node.SuccTaken != nil {
			succLabel := strconv.Itoa(node.SuccTaken.Label)
			err = graph.AddEdge(currLabel, succLabel, true, map[string]string{"label": `"True"`})
			if err != nil {
				panic(err)
			}
		}

		if node.SuccNotTaken != nil {
			succLabel := strconv.Itoa(node.SuccNotTaken.Label)
			var edgeLabel map[string]string
			if node.SuccTaken != nil {
				// It's a branch
				edgeLabel = map[string]string{"label": `"False"`}
			}
			err = graph.AddEdge(currLabel, succLabel, true, edgeLabel)
			if err != nil {
				panic(err)
			}
		}

		// Checks the preds were set correctly
		//for _, pred := range node.PredsNotTaken {
		//	predLabel := strconv.Itoa(pred.Label)
		//	graph.AddEdge(currLabel, predLabel, true, map[string]string{"label": `"FromNotTaken"`})
		//}
		//
		//for _, pred := range node.PredsTaken {
		//	predLabel := strconv.Itoa(pred.Label)
		//	graph.AddEdge(currLabel, predLabel, true, map[string]string{"label": `"FromTaken"`})
		//}
	})
}
