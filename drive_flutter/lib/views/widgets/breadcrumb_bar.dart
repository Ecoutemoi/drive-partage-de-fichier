import 'package:flutter/material.dart';
import '../../controllers/drive_controller.dart';

class BreadcrumbBar extends StatelessWidget {
  final List<Crumb> crumbs;
  final void Function(int index) onTap;

  const BreadcrumbBar({super.key, required this.crumbs, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      padding: const EdgeInsets.all(8),
      child: Row(
        children: [
          for (int i = 0; i < crumbs.length; i++) ...[
            TextButton(
              onPressed: () => onTap(i),
              child: Text(crumbs[i].name),
            ),
            if (i != crumbs.length - 1) const Text(' / '),
          ],
        ],
      ),
    );
  }
}
