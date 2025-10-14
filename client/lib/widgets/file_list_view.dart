import 'package:flutter/material.dart';

class FileListView extends StatelessWidget {
  final List<dynamic> files;

  const FileListView({super.key, required this.files});

  @override
  Widget build(BuildContext context) {
    if (files.isEmpty) {
      return const Center(
        child: Text(
          'No files found. Connect to your server.',
          style: TextStyle(color: Colors.white54),
        ),
      );
    }

    return ListView.builder(
      itemCount: files.length,
      itemBuilder: (context, index) {
        final f = files[index];
        final isFolder = f['type'] == 'folder';

        return Card(
          child: ListTile(
            leading: Icon(
              isFolder ? Icons.folder : Icons.insert_drive_file,
              color: isFolder ? Colors.amber : Colors.blue,
            ),
            title: Text(
              f['name'],
              style: const TextStyle(color: Colors.white),
            ),
            subtitle: Text(
              'Size: ${f['size']} bytes',
              style: const TextStyle(color: Colors.white54),
            ),
            onTap: isFolder
                ? () {
                    // TO DO: add folder navigation
                  }
                : null,
          ),
        );
      },
    );
  }
}