import 'package:flutter/material.dart';
import '../core/vpn_channel.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  bool _isConnected = false;
  bool _isConnecting = false;

  Future<void> _toggleConnection() async {
    setState(() {
      _isConnecting = true;
    });

    if (_isConnected) {
      bool success = await VpnChannel.stopVpn();
      if (success) {
        setState(() {
          _isConnected = false;
        });
      }
    } else {
      bool success = await VpnChannel.startVpn();
      if (success) {
        setState(() {
          _isConnected = true;
        });
      }
    }

    setState(() {
      _isConnecting = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('NEXUS CORE', style: TextStyle(fontWeight: FontWeight.bold, letterSpacing: 2)),
        centerTitle: true,
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            GestureDetector(
              onTap: _isConnecting ? null : _toggleConnection,
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 300),
                width: 200,
                height: 200,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: _isConnected ? Colors.tealAccent.withOpacity(0.2) : Colors.redAccent.withOpacity(0.2),
                  border: Border.all(
                    color: _isConnected ? Colors.tealAccent : Colors.redAccent,
                    width: 4,
                  ),
                  boxShadow: [
                    BoxShadow(
                      color: _isConnected ? Colors.tealAccent.withOpacity(0.5) : Colors.redAccent.withOpacity(0.5),
                      blurRadius: 30,
                      spreadRadius: 5,
                    )
                  ],
                ),
                child: Center(
                  child: _isConnecting
                      ? const CircularProgressIndicator(color: Colors.white)
                      : Icon(
                          Icons.power_settings_new,
                          size: 80,
                          color: _isConnected ? Colors.tealAccent : Colors.redAccent,
                        ),
                ),
              ),
            ),
            const SizedBox(height: 40),
            Text(
              _isConnected ? 'SYSTEM ONLINE' : 'SYSTEM OFFLINE',
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: _isConnected ? Colors.tealAccent : Colors.redAccent,
                letterSpacing: 3,
              ),
            ),
            const SizedBox(height: 10),
            Text(
              _isConnected ? 'Zero-Trace Routing Active' : 'Awaiting Initialization',
              style: const TextStyle(color: Colors.grey, fontSize: 16),
            ),
          ],
        ),
      ),
    );
  }
}
