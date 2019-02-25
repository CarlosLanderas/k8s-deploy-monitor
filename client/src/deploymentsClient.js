
function deploymentsClient(onMessage) {
    
    var ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = (event) => {        
        onMessage(JSON.parse(event.data))       
    };

    ws.onclose = (event) => {
      console.log(event);
      alert("Disconnected from deployments socket");
    };
}

export default deploymentsClient;
