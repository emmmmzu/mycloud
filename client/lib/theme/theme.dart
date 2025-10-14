import 'package:flutter/material.dart';

/// MyCloudApp theme configuration
final ThemeData myCloudTheme = ThemeData(
  brightness: Brightness.dark,
  scaffoldBackgroundColor: const Color(0xFF0D1B1E),
  primaryColor: const Color(0xFF00BFA5), // turquoise
  colorScheme: const ColorScheme.dark(
    primary: Color(0xFF00BFA5), // turquoise
    secondary: Color(0xFFFF5252), // red accent
    surface: Color(0xFF102C34),
    onPrimary: Colors.white,
    onSurface: Colors.white70,
    error: Color(0xFFFF1744),
  ),
  appBarTheme: const AppBarTheme(
    backgroundColor: Color(0xFF102C34),
    foregroundColor: Colors.white,
    elevation: 2,
  ),
  cardTheme: const CardThemeData(
    color: Color(0xFF102C34),
    elevation: 2,
    margin: EdgeInsets.symmetric(vertical: 6, horizontal: 4),
    shape: RoundedRectangleBorder(
      borderRadius: BorderRadius.all(Radius.circular(12)),
    ),
  ),
  elevatedButtonTheme: ElevatedButtonThemeData(
    style: ElevatedButton.styleFrom(
      backgroundColor: const Color(0xFF00BFA5),
      foregroundColor: Colors.black,
      textStyle: const TextStyle(fontWeight: FontWeight.bold),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
    ),
  ),
  floatingActionButtonTheme: const FloatingActionButtonThemeData(
    backgroundColor: Color(0xFFFF5252),
    foregroundColor: Colors.white,
  ),
  textTheme: const TextTheme(
    bodyMedium: TextStyle(color: Colors.white70, fontSize: 16),
    titleLarge: TextStyle(
      color: Color(0xFF00E5C0),
      fontWeight: FontWeight.bold,
      letterSpacing: 1.2,
    ),
  ),
);