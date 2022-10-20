# file-server-go
Create a server to allow files to be transferred between clients that connect to the server using the TCP protocol

## Protocol design

- start the server: "go run server.go"
![texto cualquiera por si no carga la imagen](https://agrotech.site/images/server.png)

- subscribe client to channel 1 to receive files: “go run client.go receive -channel 1”
![texto cualquiera por si no carga la imagen](https://agrotech.site/images/clientReceive.png)

- subscribe client to channel 1 to send files: "go run client.go send -channel 1"
![texto cualquiera por si no carga la imagen](https://agrotech.site/images/client.png)

- To send a file, write the name of the file with the path if necessary depending on the location and then enter: "miarchivo.png"
![texto cualquiera por si no carga la imagen](https://agrotech.site/images/client-send.png)

- Client received the file
![texto cualquiera por si no carga la imagen](https://agrotech.site/images/client-receive.png)

