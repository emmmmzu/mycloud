import 'package:flutter/material.dart';

/// A gentle turquoise ripple that expands outward from the center.
class AnimatedBackground extends StatefulWidget {
  const AnimatedBackground({super.key});

  @override
  State<AnimatedBackground> createState() => _AnimatedBackgroundState();
}

class _AnimatedBackgroundState extends State<AnimatedBackground>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 8),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, child) {
        return CustomPaint(
          painter: _RipplePainter(_controller.value),
          child: Container(),
        );
      },
    );
  }
}

class _RipplePainter extends CustomPainter {
  final double progress;

  _RipplePainter(this.progress);

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final maxRadius = size.width * 1.5; // how far the ripple expands

    // Background
    final Paint backgroundPaint = Paint()
      ..color = const Color(0xFF0D1B1E); // dark teal background
    canvas.drawRect(Rect.fromLTWH(0, 0, size.width, size.height), backgroundPaint);

    // Create a soft turquoise ripple effect
    final rippleCount = 20;
    for (int i = 0; i < rippleCount; i++) {
      final double rippleProgress = (progress + i / rippleCount) % 1;
      final double radius = rippleProgress * maxRadius;

      final Paint ripplePaint = Paint()
        ..shader = RadialGradient(
          colors: [
            const Color(0xFF0D1B1E),
            const Color(0xFF00E5C0),
            const Color(0xFF0D1B1E),
          ],
          stops: const [0.8, 0.9, 1],
        ).createShader(Rect.fromCircle(center: center, radius: radius));

      canvas.drawCircle(center, radius, ripplePaint);
    }
  }

  @override
  bool shouldRepaint(covariant _RipplePainter oldDelegate) => true;
}
