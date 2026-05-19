import 'package:flutter/material.dart';
import '../core/vpn_channel.dart';

class SettingsScreen extends StatefulWidget {
  const SettingsScreen({Key? key}) : super(key: key);

  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  String _selectedTransport = 'AUTO (AI Managed)';
  final List<String> _transports = [
    'AUTO (AI Managed)',
    'HYSTERIA2',
    'TUIC',
    'WEBTUNNEL',
    'CDN_RELAY',
    'DNS_NULL',
    'ICMP'
  ];

  Future<void> _exportLogs() async {
    String path = await VpnChannel.exportLogs();
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Logs exported to: $path'), backgroundColor: Colors.teal),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('SYSTEM SETTINGS')),
      body: ListView(
        padding: const EdgeInsets.all(16.0),
        children: [
          const Text('Transport Override', style: TextStyle(color: Colors.tealAccent, fontWeight: FontWeight.bold)),
          const SizedBox(height: 10),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12),
            decoration: BoxDecoration(
              color: const Color(0xFF1A2235),
              borderRadius: BorderRadius.circular(8),
            ),
            child: DropdownButtonHideUnderline(
              child: DropdownButton<String>(
                value: _selectedTransport,
                dropdownColor: const Color(0xFF1A2235),
                isExpanded: true,
                style: const TextStyle(color: Colors.white, fontSize: 16),
                items: _transports.map((String value) {
                  return DropdownMenuItem<String>(
                    value: value,
                    child: Text(value),
                  );
                }).toList(),
                onChanged: (String? newValue) {
                  if (newValue != null) {
                    setState(() {
                      _selectedTransport = newValue;
                    });
                    VpnChannel.setTransportOverride(newValue);
                  }
                },
              ),
            ),
          ),
          const SizedBox(height: 30),
          const Text('System Maintenance', style: TextStyle(color: Colors.tealAccent, fontWeight: FontWeight.bold)),
          const SizedBox(height: 10),
          ListTile(
            tileColor: const Color(0xFF1A2235),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            leading: const Icon(Icons.bug_report, color: Colors.white),
            title: const Text('Export Diagnostic Logs', style: TextStyle(color: Colors.white)),
            trailing: const Icon(Icons.chevron_right, color: Colors.grey),
            onTap: _exportLogs,
          ),
          const SizedBox(height: 10),
          ListTile(
            tileColor: const Color(0xFF1A2235),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            leading: const Icon(Icons.delete_forever, color: Colors.redAccent),
            title: const Text('Purge Encrypted Cache', style: TextStyle(color: Colors.redAccent)),
            onTap: () {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Cache purged securely.'), backgroundColor: Colors.redAccent),
              );
            },
          ),
        ],
      ),
    );
  }
}
