import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'core/config.dart';
import 'controllers/drive_controller.dart';
import 'services/drive_api.dart';
import 'views/drive_screen.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    final api = DriveApi(AppConfig.apiBaseUrl);

    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => DriveController(api: api),
        ),
      ],
      child: MaterialApp(
        debugShowCheckedModeBanner: false,
        title: 'Drive MVC',
        theme: ThemeData(
          useMaterial3: true,
          colorSchemeSeed: Colors.blue,
        ),
        home: const DriveScreen(),
      ),
    );
  }
}
