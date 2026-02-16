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
