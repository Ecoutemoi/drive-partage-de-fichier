import 'dart:convert';
import 'package:http/http.dart' as http;

class FolderItem {
  final int id;
  final int userId;
  final String name;
  final int? parentId;

  FolderItem({
    required this.id,
    required this.userId,
    required this.name,
    required this.parentId,
  });

  factory FolderItem.fromJson(Map<String, dynamic> j) => FolderItem(
        id: (j['id'] as num).toInt(),
        userId: (j['user_id'] as num).toInt(),
        name: (j['name'] ?? '').toString(),
        parentId: j['parent_id'] == null ? null : (j['parent_id'] as num).toInt(),
      );
}

class FileItem {
  final int id;
  final int userId;
  final int? folderId;
  final String originalName;
  final String mimeType;
  final int sizeBytes;

  FileItem({
    required this.id,
    required this.userId,
    required this.folderId,
    required this.originalName,
    required this.mimeType,
    required this.sizeBytes,
  });

  factory FileItem.fromJson(Map<String, dynamic> j) => FileItem(
        id: (j['id'] as num).toInt(),
        userId: (j['user_id'] as num).toInt(),
        folderId: j['folder_id'] == null ? null : (j['folder_id'] as num).toInt(),
        originalName: (j['original_name'] ?? '').toString(),
        mimeType: (j['mime_type'] ?? '').toString(),
        sizeBytes: (j['size_bytes'] as num).toInt(),
      );
}

class DriveApi {
  final String baseUrl;
  DriveApi(this.baseUrl);

  Uri _u(String path, [Map<String, String>? q]) => Uri.parse('$baseUrl$path').replace(queryParameters: q);

  Future<Map<String, dynamic>> listFolder({int? parentId}) async {
    final q = <String, String>{};
    if (parentId != null) q['parent_id'] = parentId.toString();

    final res = await http.get(_u('/folders/list', q));
    if (res.statusCode != 200) {
      throw Exception('List error ${res.statusCode}: ${res.body}');
    }
    return jsonDecode(res.body) as Map<String, dynamic>;
  }

  Future<void> createFolder({required String name, int? parentId}) async {
    final body = {
      'name': name,
      'parent_id': parentId,
    };
    final res = await http.post(
      _u('/Createfolders'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode(body),
    );
    if (res.statusCode != 201) {
      throw Exception('Create folder error ${res.statusCode}: ${res.body}');
    }
  }

  Future<void> renameFolder({required int folderId, required String newName}) async {
    final res = await http.put(
      _u('/folders/rename'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'doss_id': folderId, 'new_name': newName}),
    );
    if (res.statusCode != 200) {
      throw Exception('Rename folder error ${res.statusCode}: ${res.body}');
    }
  }

  Future<void> moveFolder({required int folderId, int? parentId}) async {
    final res = await http.put(
      _u('/folders/move'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'folder_id': folderId, 'parent_id': parentId}),
    );
    if (res.statusCode != 200) {
      throw Exception('Move folder error ${res.statusCode}: ${res.body}');
    }
  }

  Future<void> deleteFolder({required int folderId}) async {
    // ton API exige users_id -> vu que pas de login, on met 1
    final res = await http.delete(
      _u('/folders/delete'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'folder_id': folderId, 'users_id': 1}),
    );
    if (res.statusCode != 200) {
      throw Exception('Delete folder error ${res.statusCode}: ${res.body}');
    }
  }

  Future<void> renameFile({required int fileId, required String newName}) async {
    final res = await http.put(
      _u('/files/rename'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'file_id': fileId, 'new_name': newName}),
    );
    if (res.statusCode != 200) {
      throw Exception('Rename file error ${res.statusCode}: ${res.body}');
    }
  }

  Future<void> moveFile({required int fileId, int? folderId}) async {
    final res = await http.put(
      _u('/files/move'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'file_id': fileId, 'folder_id': folderId}),
    );
    if (res.statusCode != 200) {
      throw Exception('Move file error ${res.statusCode}: ${res.body}');
    }
  }

  Future<void> deleteFile({required int fileId}) async {
    final res = await http.delete(
      _u('/files/delete'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'file_id': fileId, 'user_id': 1}),
    );
    if (res.statusCode != 200) {
      throw Exception('Delete file error ${res.statusCode}: ${res.body}');
    }
  }

  Future<String> createShareLink({required int fileId}) async {
    // IMPORTANT: ton Go utilise /files/{id}/share
    final res = await http.post(_u('/files/$fileId/share'));
    if (res.statusCode != 201) {
      throw Exception('Share error ${res.statusCode}: ${res.body}');
    }
    final data = jsonDecode(res.body) as Map<String, dynamic>;
    return (data['url'] ?? '').toString();
  }
}
