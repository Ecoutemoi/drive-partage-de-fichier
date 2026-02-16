import 'package:flutter/material.dart';
import '../../core/utils.dart';
import '../../models/file_item.dart';

class FileTile extends StatelessWidget {
  final FileItem file;
  final VoidCallback onDownload;
  final VoidCallback onShare;
  final VoidCallback onMove; // ✅ nouveau
  final VoidCallback onRename;
  final VoidCallback onDelete;

  const FileTile({
    super.key,
    required this.file,
    required this.onDownload,
    required this.onShare,
    required this.onMove, // ✅ nouveau
    required this.onRename,
    required this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: const Icon(Icons.insert_drive_file),
      title: Text(file.originalName),
      subtitle: Text('${file.mimeType} • ${fmtSize(file.sizeBytes)}'),
      onTap: onDownload,
      trailing: PopupMenuButton<String>(
        onSelected: (v) {
          if (v == 'download') onDownload();
          if (v == 'share') onShare();
          if (v == 'move') onMove(); // ✅ nouveau
          if (v == 'rename') onRename();
          if (v == 'delete') onDelete();
        },
        itemBuilder: (_) => const [
          PopupMenuItem(value: 'download', child: Text('Télécharger')),
          PopupMenuItem(value: 'share', child: Text('Partager')),
          PopupMenuItem(value: 'move', child: Text('Déplacer')), // ✅ nouveau
          PopupMenuItem(value: 'rename', child: Text('Renommer')),
          PopupMenuItem(value: 'delete', child: Text('Supprimer')),
        ],
      ),
    );
  }
}
