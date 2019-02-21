
function deploymentsClient(onMessage) {
    
    var ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = (event) => {        
        onMessage(JSON.parse(event.data))       
    }
}


export default deploymentsClient;
