import 'package:flutter/material.dart';
import '../../models/folder_item.dart';

class FolderTile extends StatelessWidget {
  final FolderItem folder;
  final VoidCallback onOpen;
  final VoidCallback onRename;
  final VoidCallback onMove; // ✅ nouveau
  final VoidCallback onDelete;

  const FolderTile({
    super.key,
    required this.folder,
    required this.onOpen,
    required this.onRename,
    required this.onMove,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: const Icon(Icons.folder),
      title: Text(folder.name),
      onTap: onOpen,
      trailing: PopupMenuButton<String>(
        onSelected: (v) {
          if (v == 'rename') onRename();
          if (v == 'move') onMove();
          if (v == 'delete') onDelete();
        },
        itemBuilder: (_) => const [
          PopupMenuItem(value: 'rename', child: Text('Renommer')),
          PopupMenuItem(value: 'move', child: Text('Déplacer')),
          PopupMenuItem(value: 'delete', child: Text('Supprimer')),
        ],
      ),
    );
  }
}
