import 'dart:async';
import 'package:flutter/material.dart';
import '../core/vpn_channel.dart';

class DiagnosticsDashboard extends StatefulWidget {
  const DiagnosticsDashboard({Key? key}) : super(key: key);

  @override
  State<DiagnosticsDashboard> createState() => _DiagnosticsDashboardState();
}

class _DiagnosticsDashboardState extends State<DiagnosticsDashboard> {
  Map<String, dynamic> _diagnostics = {};
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _fetchDiagnostics();
    _timer = Timer.periodic(const Duration(seconds: 2), (timer) {
      _fetchDiagnostics();
    });
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  Future<void> _fetchDiagnostics() async {
    final data = await VpnChannel.getDiagnostics();
    if (mounted) {
      setState(() {
        _diagnostics = data;
      });
    }
  }

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
    return Card(
      color: const Color(0xFF1A2235),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            Icon(icon, size: 40, color: color),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: const TextStyle(color: Colors.grey, fontSize: 14)),
                  const SizedBox(height: 4),
                  Text(value, style: const TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('LIVE DIAGNOSTICS')),
      body: ListView(
        padding: const EdgeInsets.all(16.0),
        children: [
          _buildStatCard(
            'Network Health (AI Q-Value)',
            _diagnostics['q_value']?.toString() ?? '0.0',
            Icons.health_and_safety,
            Colors.greenAccent,
          ),
          const SizedBox(height: 12),
          _buildStatCard(
            'Active Transport Level',
            _diagnostics['active_transport'] ?? 'UNKNOWN',
            Icons.route,
            Colors.blueAccent,
          ),
          const SizedBox(height: 12),
          _buildStatCard(
            'Current Latency',
            '${_diagnostics['latency_ms'] ?? 0} ms',
            Icons.timer,
            Colors.orangeAccent,
          ),
          const SizedBox(height: 12),
          _buildStatCard(
            'Tor Bridge Status',
            _diagnostics['bridge_status'] ?? 'OFFLINE',
            Icons.security,
            Colors.purpleAccent,
          ),
        ],
      ),
    );
  }
}
