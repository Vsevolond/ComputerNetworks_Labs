// ws://echo.websocket.org - не рабочий
// wss://echo.websocket.events
// ws://151.248.113.144:8765
// flutter run

import 'package:flutter/material.dart';
import 'package:web_socket_channel/io.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: MyHomePage(),
    );
  }
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  //WebSocketChannel channel = IOWebSocketChannel.connect('wss://echo.websocket.events');
  //WebSocketChannel channel = IOWebSocketChannel.connect('ws://151.248.113.144:8765');
  WebSocketChannel channel = IOWebSocketChannel.connect('ws://151.248.113.144:8000');
  TextEditingController controller = TextEditingController();
  final List<String> list = [];

  @override
  void initState() {
    print("initState...");
    super.initState();
    channel.stream.listen((data) => setState(() => list.add(data)));
    print(list);
    print("started to listen");
  }

  void sendData() {
    if (controller.text.isNotEmpty) {
      channel.sink.add(controller.text);
      print(list);
      print("-------->"+controller.text);
      controller.text = "";
    }
  }

  @override
  void dispose() {
    print("closeing connection...");
    channel.sink.close();
    print("done");
   super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('WebSocket Example'),
      ),
      body: Container(
        padding: EdgeInsets.all(20.0),
        child: Column(
          children: <Widget>[
            Form(
              child: TextFormField(
                controller: controller,
                decoration: InputDecoration(
                  labelText: "Send to WebSocket",
                ),
              ),
            ),
            Column(
              children: list.map((data) => Text(data)).toList(),
            )
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        child: Icon(Icons.send),
        onPressed: () {
          sendData();
        },
      ),
    );
  }
}
  
