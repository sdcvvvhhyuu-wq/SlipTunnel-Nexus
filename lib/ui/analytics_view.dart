import 'package:flutter/material.dart';

class AnalyticsView extends StatelessWidget {
  const AnalyticsView({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('TRAFFIC ANALYTICS')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Session Data Usage', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white)),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                _buildDataCircle('TX', '142 MB', Colors.tealAccent),
                _buildDataCircle('RX', '890 MB', Colors.blueAccent),
              ],
            ),
            const SizedBox(height: 40),
            const Text('Transport Fallback History', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Colors.white)),
            const SizedBox(height: 10),
            Expanded(
              child: ListView(
                children: const [
                  ListTile(
                    leading: Icon(Icons.arrow_downward, color: Colors.redAccent),
                    title: Text('Escalated to CDN_RELAY', style: TextStyle(color: Colors.white)),
                    subtitle: Text('DPI Interference Detected - 10:42 AM', style: TextStyle(color: Colors.grey)),
                  ),
                  ListTile(
                    leading: Icon(Icons.arrow_upward, color: Colors.greenAccent),
                    title: Text('Restored to HYSTERIA2', style: TextStyle(color: Colors.white)),
                    subtitle: Text('Network Stabilized - 11:15 AM', style: TextStyle(color: Colors.grey)),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDataCircle(String label, String value, Color color) {
    return Container(
      width: 120,
      height: 120,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        border: Border.all(color: color, width: 4),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(label, style: TextStyle(color: color, fontSize: 20, fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Text(value, style: const TextStyle(color: Colors.white, fontSize: 16)),
        ],
      ),
    );
  }
}
