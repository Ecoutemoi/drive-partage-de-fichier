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
