package funcs

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
	"gat/database"
	"fmt"
)

func home(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message": "Lol funny api",
	})
	return
}

func addNode(c *gin.Context){
	var node db2.Node

	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error in binding": err.Error(),})
		return
	}

	fmt.Println(node)

	res, err := db2.AddNode(&node)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error in database:": err.Error(),})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"node added": res})
	return
}

func allNodes(c *gin.Context){
	nodes, err := db2.AllNodes()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
	return
}

func addEdge(c *gin.Context){

	type Body struct {
		Nodes 	[]db2.Node 	`jason:"nodes"`
		Dist  	int		`jason:"dist"`
	}

	var body Body

	if err := c.ShouldBindJSON(&body); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	if len(body.Nodes) != 2{
		c.JSON(http.StatusBadRequest, gin.H{"error": "expecting an array of 2 user objects"})
        return
	}

	node1 := body.Nodes[0] // start and end node
	node2 := body.Nodes[1]

	fmt.Println("one-", node1, "two-",node2)
	nodes, err := db2.AddEdge(&node1, &node2)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nodes" : nodes,})
	return
}

func SpreadRadius(c * gin.Context){
	start_id, err := strconv.Atoi(c.Param("start"))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}
	
	start, err := db2.FindNode(int64(start_id))
	
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	mode := c.Param("mode") // debug to turn multithreading on and off

	var paths [][]*db2.Node

	if mode == "0"{
		paths = db2.SpreadRadiusSingle(start, 3, 0, make([]*db2.Node, 0), make([][]*db2.Node, 0))	
	} else{
		paths = db2.SpreadRadius(start, 3, 0, make([]*db2.Node, 0), make([][]*db2.Node, 0))
	}
	
	readable_paths := make([][]int64, 0)

	for _, e := range paths{
		readable_paths = append(readable_paths, db2.IdSliceFromNodeSlice(e))
	}

	c.JSON(http.StatusOK, gin.H{"paths":readable_paths,})
	return
}

func connect(c * gin.Context){
	n1, err1 := strconv.Atoi(c.Param("n1"))
	n2, err2 := strconv.Atoi(c.Param("n2"))

	if err1 != nil|| err2 != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error(),})
		return

	}

	err := db2.ConnectNodes(int64(n1),int64(n2))
	
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node1" : n1, "node2": n2,})
	return
}

func SetupRouter() *gin.Engine{
	r := gin.Default()
	r.GET("/", home)
	r.POST("/addNode", addNode)
	r.POST("/addEdge", addEdge)
	r.GET("/allNodes", allNodes)
	r.PUT("/addEdge/:n1/:n2", connect)
	r.GET("/spreadRadius/:start/:mode", SpreadRadius)
	return r
}
