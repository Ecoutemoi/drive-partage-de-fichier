import 'package:file_picker/file_picker.dart';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';

import '../models/folder_item.dart';
import '../models/file_item.dart';
import '../services/drive_api.dart';

class Crumb {
  final String name;
  final int? id;
  const Crumb({required this.name, required this.id});
}

class DriveController extends ChangeNotifier {
  final DriveApi api;

  DriveController({required this.api});

  bool loading = false;
  int? currentParentId;
  List<Crumb> crumbs = const [Crumb(name: 'Racine', id: null)];

  List<FolderItem> folders = [];
  List<FileItem> files = [];

  Future<void> refresh() async {
    loading = true;
    notifyListeners();
    try {
      final data = await api.listFolder(parentId: currentParentId);
      folders = (data['folders'] as List? ?? [])
          .map((e) => FolderItem.fromJson(e as Map<String, dynamic>))
          .toList();
      files = (data['files'] as List? ?? [])
          .map((e) => FileItem.fromJson(e as Map<String, dynamic>))
          .toList();
    } finally {
      loading = false;
      notifyListeners();
    }
  }

  void openFolder(FolderItem f) {
    currentParentId = f.id;
    crumbs = [...crumbs, Crumb(name: f.name, id: f.id)];
    notifyListeners();
    refresh();
  }

  void goToCrumb(int index) {
    final c = crumbs[index];
    currentParentId = c.id;
    crumbs = crumbs.sublist(0, index + 1);
    notifyListeners();
    refresh();
  }

  Future<void> createFolder(String name) async {
    await api.createFolder(name: name, parentId: currentParentId);
    await refresh();
  }

  Future<void> renameFolder(int folderId, String newName) async {
    await api.renameFolder(folderId: folderId, newName: newName);
    await refresh();
  }

  Future<void> deleteFolder(int folderId) async {
    await api.deleteFolder(folderId: folderId);
    await refresh();
  }

  Future<void> renameFile(int fileId, String newName) async {
    await api.renameFile(fileId: fileId, newName: newName);
    await refresh();
  }

  Future<void> deleteFile(int fileId) async {
    await api.deleteFile(fileId: fileId);
    await refresh();
  }

  Future<String> shareFile(int fileId) async {
    return api.createShareLink(fileId: fileId);
  }

  Future<void> downloadFile(int fileId) async {
    final url = Uri.parse('${api.baseUrl}/files/$fileId/download');
    await launchUrl(url, mode: LaunchMode.externalApplication);
  }

  Future<void> uploadPickedFile() async {
    final picked = await FilePicker.platform.pickFiles(withData: true);
    if (picked == null || picked.files.isEmpty) return;

    final file = picked.files.first;
    if (file.bytes == null) {
      throw Exception('Upload impossible: bytes null');
    }

    final uri = Uri.parse('${api.baseUrl}/upload');
    final req = http.MultipartRequest('POST', uri);

    if (currentParentId != null) {
      req.fields['folder_id'] = currentParentId.toString();
    }

    req.files.add(
      http.MultipartFile.fromBytes('file', file.bytes!, filename: file.name),
    );

    final res = await req.send();
    final body = await res.stream.bytesToString();

    if (res.statusCode != 201) {
      throw Exception('Erreur upload ${res.statusCode}: $body');
    }

    await refresh();
  }

  /// ✅ Déplacer un fichier (folderId null => racine)
  Future<void> moveFileToFolder(int fileId, int? folderId) async {
    await api.moveFile(fileId: fileId, folderId: folderId);
    await refresh();
  }

  /// ✅ Déplacer un dossier (parentId null => racine)
  Future<void> moveFolderToParent(int folderId, int? parentId) async {
    await api.moveFolder(folderId: folderId, parentId: parentId);
    await refresh();
  }
}
