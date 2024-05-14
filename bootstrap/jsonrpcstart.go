package bootstrap


import (
	go_jsonrpc "github.com/iloveswift/go-jsonrpc"
	"testproject/app/services"
	"testproject/global"
	//"fmt"
   )
   

   
   func JsonRpcConnect(){
	s, _ := go_jsonrpc.NewServer("http", global.App.Config.JsonRpc.Host, "3345")
	s.Register(new(services.Go))
	go func() {  
		s.Start()
	}() 

	
   }