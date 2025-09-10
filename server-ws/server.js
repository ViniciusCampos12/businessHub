const WebSocket = require('ws');

const wss = new WebSocket.Server({ port: 8080 });

const clients = [];

console.log('Servidor WebSocket iniciado na porta 8080');

function broadcast(message, excludeClient = null) {
    const messageStr = typeof message === 'string' ? message : JSON.stringify(message);
    
    clients.forEach((client, index) => {
        if (client.readyState === WebSocket.OPEN && client !== excludeClient) {
            try {
                client.send(messageStr);
            } catch (error) {
                console.error(`Erro ao enviar mensagem para cliente ${index}:`, error);
                const clientIndex = clients.indexOf(client);
                if (clientIndex > -1) {
                    clients.splice(clientIndex, 1);
                }
            }
        }
    });
}

wss.on('connection', (ws) => {
    console.log(`Nova conexão estabelecida. Total de clientes: ${clients.length + 1}`);
    
    clients.push(ws);
    
    ws.send(JSON.stringify({
        type: 'welcome',
        message: 'Conectado ao servidor WebSocket',
        clientCount: clients.length,
        timestamp: new Date().toISOString()
    }));
    
    broadcast({
        type: 'client_joined',
        message: 'Um novo cliente se conectou',
        clientCount: clients.length,
        timestamp: new Date().toISOString()
    }, ws);
    
    ws.on('message', (data) => {
        let message;
        try {
            message = JSON.parse(data.toString());
        } catch (error) {
            message = {
                type: 'message',
                content: data.toString(),
                timestamp: new Date().toISOString()
            };
        }
        
        console.log('Mensagem recebida:', message);
        
        broadcast({
            type: 'broadcast',
            data: message,
            timestamp: new Date().toISOString()
        }, ws);
    });
    
    ws.on('close', () => {
        const clientIndex = clients.indexOf(ws);
        if (clientIndex > -1) {
            clients.splice(clientIndex, 1);
        }
        
        console.log(`Cliente desconectado. Total de clientes: ${clients.length}`);

        broadcast({
            type: 'client_left',
            message: 'Um cliente se desconectou',
            clientCount: clients.length,
            timestamp: new Date().toISOString()
        });
    });
    
    ws.on('error', (error) => {
        console.error('Erro na conexão WebSocket:', error);
        const clientIndex = clients.indexOf(ws);
        if (clientIndex > -1) {
            clients.splice(clientIndex, 1);
        }
    });
    
    ws.on('pong', () => {
        ws.isAlive = true;
    });
});

const interval = setInterval(() => {
    clients.forEach((ws, index) => {
        if (ws.isAlive === false) {
            console.log(`Removendo conexão inativa ${index}`);
            const clientIndex = clients.indexOf(ws);
            if (clientIndex > -1) {
                clients.splice(clientIndex, 1);
            }
            return ws.terminate();
        }
        
        ws.isAlive = false;
        ws.ping();
    });
}, 30000);

wss.on('close', () => {
    clearInterval(interval);
});

setInterval(() => {
    console.log(`Status: ${clients.length} clientes conectados`);
    
    if (clients.length > 0) {
        broadcast({
            type: 'status',
            message: 'Status do servidor',
            clientCount: clients.length,
            uptime: process.uptime(),
            timestamp: new Date().toISOString()
        });
    }
}, 60000);