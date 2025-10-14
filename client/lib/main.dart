import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'theme/theme.dart';
import 'widgets/animated_background.dart';
import 'widgets/server_connection_card.dart';
import 'widgets/file_list_view.dart';

void main() {
  runApp(MyCloudApp());
}

class MyCloudApp extends StatelessWidget {
  const MyCloudApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MyCloud Client',
      theme: myCloudTheme,
      home: ServerPage(),
    );
  }
}

class ServerPage extends StatefulWidget {
  const ServerPage({super.key});

  @override
  // ignore: library_private_types_in_public_api
  _ServerPageState createState() => _ServerPageState();
}

class _ServerPageState extends State<ServerPage> {
  final TextEditingController _urlController = TextEditingController();
  List<dynamic> files = [];
  bool loading = false;
  bool buttonError = false;

  @override
  void initState() {
    super.initState();
    _loadServerUrl();
  }

  Future<void> _loadServerUrl() async {
    final prefs = await SharedPreferences.getInstance();
    _urlController.text = prefs.getString('serverUrl') ?? '';
  }

  Future<void> _saveServerUrl() async {
    final prefs = await SharedPreferences.getInstance();
    prefs.setString('serverUrl', _urlController.text);
  }

  Future<void> fetchFiles() async {
    final url = '${_urlController.text}/list?path=/';
    setState(() => loading = true);

    try {
      final response = await http.get(Uri.parse(url));
      if (response.statusCode == 200) {
        setState(() {
          files = jsonDecode(response.body);
        });
      } else {
        _showError('Server returned ${response.statusCode}');
        _flashError();
      }
    } catch (e) {
      _showError('Failed to connect: $e');
      _flashError();
    }

    setState(() => loading = false);
  }

  void _flashError() {
    // temporary red pulse effect
    setState(() => buttonError = true);
    Future.delayed(const Duration(milliseconds: 600), () {
      setState(() => buttonError = false);
    });
  }

  void _showError(String msg) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(msg)));
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        const AnimatedBackground(),
        Scaffold(
          backgroundColor: Colors.transparent,
          appBar: AppBar(
            title: const Text(
              'MyCloud',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.w600,
                letterSpacing: 1.2,
              ),
            ),
            centerTitle: true,
          ),
          floatingActionButton: FloatingActionButton(
            onPressed: () {
              // TO DO: implement upload picker
            },
            child: const Icon(Icons.add),
          ),
          body: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                ServerConnectionCard(
                  urlController: _urlController,
                  loading: loading,
                  error: buttonError,
                  onConnect: () {
                    _saveServerUrl();
                    fetchFiles();
                  },
                ),
                const SizedBox(height: 12),
                Expanded(child: FileListView(files: files)),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
