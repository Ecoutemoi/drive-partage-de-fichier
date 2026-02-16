import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../controllers/drive_controller.dart';
import 'widgets/breadcrumb_bar.dart';
import 'widgets/file_tile.dart';
import 'widgets/folder_tile.dart';

class DriveScreen extends StatefulWidget {
  const DriveScreen({super.key});

  @override
  State<DriveScreen> createState() => _DriveScreenState();
}

class _DriveScreenState extends State<DriveScreen> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => context.read<DriveController>().refresh());
  }

  void _snack(String s) =>
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(s)));

  Future<String?> _promptText(String title, {String initial = ''}) async {
    final ctrl = TextEditingController(text: initial);
    return showDialog<String>(
      context: context,
      builder: (_) => AlertDialog(
        title: Text(title),
        content: TextField(controller: ctrl, autofocus: true),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Annuler')),
          FilledButton(
              onPressed: () => Navigator.pop(context, ctrl.text.trim()),
              child: const Text('OK')),
        ],
      ),
    );
  }

  Future<bool?> _confirm(String title, String content) {
    return showDialog<bool>(
      context: context,
      builder: (_) => AlertDialog(
        title: Text(title),
        content: Text(content),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text('Annuler')),
          FilledButton(
              onPressed: () => Navigator.pop(context, true),
              child: const Text('Confirmer')),
        ],
      ),
    );
  }

  /// ✅ Choisir destination dossier pour un FILE
  /// null => annuler ; -1 => racine ; sinon => id dossier
  Future<int?> _pickFolderDestinationForFile() {
    final c = context.read<DriveController>();

    return showDialog<int>(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Déplacer le fichier vers...'),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView(
            shrinkWrap: true,
            children: [
              ListTile(
                leading: const Icon(Icons.home),
                title: const Text('Racine'),
                onTap: () => Navigator.pop(context, -1),
              ),
              const Divider(),
              for (final f in c.folders)
                ListTile(
                  leading: const Icon(Icons.folder),
                  title: Text(f.name),
                  onTap: () => Navigator.pop(context, f.id),
                ),
            ],
          ),
        ),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Annuler')),
        ],
      ),
    );
  }

  /// ✅ Choisir destination parent pour un DOSSIER
  /// On évite de proposer le dossier lui-même
  Future<int?> _pickParentDestinationForFolder({required int movingFolderId}) {
    final c = context.read<DriveController>();

    return showDialog<int>(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('Déplacer le dossier vers...'),
        content: SizedBox(
          width: double.maxFinite,
          child: ListView(
            shrinkWrap: true,
            children: [
              ListTile(
                leading: const Icon(Icons.home),
                title: const Text('Racine'),
                onTap: () => Navigator.pop(context, -1),
              ),
              const Divider(),
              for (final f in c.folders)
                if (f.id != movingFolderId)
                  ListTile(
                    leading: const Icon(Icons.folder),
                    title: Text(f.name),
                    onTap: () => Navigator.pop(context, f.id),
                  ),
            ],
          ),
        ),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('Annuler')),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final c = context.watch<DriveController>();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Drive (MVC, sans login)'),
        actions: [
          IconButton(
            onPressed: c.loading
                ? null
                : () async {
                    try {
                      await c.refresh();
                    } catch (e) {
                      _snack(e.toString());
                    }
                  },
            icon: const Icon(Icons.refresh),
          ),
          IconButton(
            onPressed: c.loading
                ? null
                : () async {
                    final name = await _promptText('Nouveau dossier');
                    if (name == null || name.isEmpty) return;
                    try {
                      await c.createFolder(name);
                    } catch (e) {
                      _snack(e.toString());
                    }
                  },
            icon: const Icon(Icons.create_new_folder),
          ),
          IconButton(
            onPressed: c.loading
                ? null
                : () async {
                    try {
                      await c.uploadPickedFile();
                    } catch (e) {
                      _snack(e.toString());
                    }
                  },
            icon: const Icon(Icons.upload),
          ),
        ],
      ),
      body: Column(
        children: [
          BreadcrumbBar(crumbs: c.crumbs, onTap: c.goToCrumb),
          const Divider(height: 1),
          Expanded(
            child: c.loading
                ? const Center(child: CircularProgressIndicator())
                : ListView(
                    children: [
                      const Padding(
                        padding: EdgeInsets.all(12),
                        child: Text(
                          'Dossiers',
                          style: TextStyle(
                              fontSize: 16, fontWeight: FontWeight.w600),
                        ),
                      ),

                      // ✅ DOSSIERS
                      for (final f in c.folders)
                        FolderTile(
                          folder: f,
                          onOpen: () => c.openFolder(f),
                          onRename: () async {
                            final name = await _promptText('Renommer dossier',
                                initial: f.name);
                            if (name == null ||
                                name.isEmpty ||
                                name == f.name) return;
                            try {
                              await c.renameFolder(f.id, name);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },

                          // ✅ DÉPLACER DOSSIER
                          onMove: () async {
                            final dest = await _pickParentDestinationForFolder(
                                movingFolderId: f.id);
                            if (dest == null) return;

                            try {
                              final parentId = (dest == -1) ? null : dest;
                              await c.moveFolderToParent(f.id, parentId);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },

                          onDelete: () async {
                            final ok = await _confirm('Supprimer dossier',
                                'Supprimer "${f.name}" ? (si non vide: conflit)');
                            if (ok != true) return;
                            try {
                              await c.deleteFolder(f.id);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },
                        ),

                      const Divider(height: 1),
                      const Padding(
                        padding: EdgeInsets.all(12),
                        child: Text(
                          'Fichiers',
                          style: TextStyle(
                              fontSize: 16, fontWeight: FontWeight.w600),
                        ),
                      ),

                      // ✅ FICHIERS
                      for (final f in c.files)
                        FileTile(
                          file: f,
                          onDownload: () async {
                            try {
                              await c.downloadFile(f.id);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },
                          onShare: () async {
                            try {
                              final url = await c.shareFile(f.id);
                              await showDialog<void>(
                                context: context,
                                builder: (_) => AlertDialog(
                                  title: const Text('Lien de partage'),
                                  content: SelectableText(url),
                                  actions: [
                                    TextButton(
                                      onPressed: () => Navigator.pop(context),
                                      child: const Text('Fermer'),
                                    ),
                                  ],
                                ),
                              );
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },

                          // ✅ DÉPLACER FICHIER
                          onMove: () async {
                            final dest = await _pickFolderDestinationForFile();
                            if (dest == null) return;

                            try {
                              final folderId = (dest == -1) ? null : dest;
                              await c.moveFileToFolder(f.id, folderId);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },

                          onRename: () async {
                            final name = await _promptText('Renommer fichier',
                                initial: f.originalName);
                            if (name == null ||
                                name.isEmpty ||
                                name == f.originalName) return;
                            try {
                              await c.renameFile(f.id, name);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },
                          onDelete: () async {
                            final ok = await _confirm('Supprimer fichier',
                                'Supprimer "${f.originalName}" ?');
                            if (ok != true) return;
                            try {
                              await c.deleteFile(f.id);
                            } catch (e) {
                              _snack(e.toString());
                            }
                          },
                        ),

                      const SizedBox(height: 24),
                    ],
                  ),
          ),
        ],
      ),
    );
  }
}
