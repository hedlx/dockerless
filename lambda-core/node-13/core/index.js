import http from 'http'

import lambda from './lambda';


const port = 3000

const requestHandler = (request, response) => {
    response.on('error', (err) => {
        console.error(err);
    });

    if (request.method === 'POST') {
        const chunks = [];

        request
            .on('error', (err) => {
                console.error(err);
            })
            .on('data', (chunk) => {
                chunks.push(chunk);
            })
            .on('end', () => {
                try {
                    const data = JSON.parse(Buffer.concat(chunks).toString());
                    const res = lambda(data);

                    response.writeHead(200, {'Content-Type': 'application/json'})
                    response.write(JSON.stringify(res));
                    response.end();
                } catch (e) {
                    response.writeHead(500, {'Content-Type': 'application/json'})
                    response.write(JSON.stringify({error: e}));
                    response.end();
                }
            });
    } else {
        response.writeHead(501)
        response.end();
    }
}

const server = http.createServer(requestHandler)

server.listen(port, err => {
  if (err) {
    return console.error(err)
  }

  console.log(`server is listening on ${port}`)
})