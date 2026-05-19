import 'dart:convert';
import 'package:flutter/services.dart';

class VpnChannel {
  static const MethodChannel _channel = MethodChannel('com.argotunnel/vpn');

  static Future<bool> startVpn() async {
    try {
      final String result = await _channel.invokeMethod('startVpn');
      return result == "VPN Started";
    } on PlatformException {
      return false;
    }
  }

  static Future<bool> stopVpn() async {
    try {
      final String result = await _channel.invokeMethod('stopVpn');
      return result == "VPN Stopped";
    } on PlatformException {
      return false;
    }
  }

  static Future<Map<String, dynamic>> getDiagnostics() async {
    try {
      final String result = await _channel.invokeMethod('getDiagnostics');
      return jsonDecode(result);
    } on PlatformException {
      return {
        "network_health": "UNKNOWN",
        "latency_ms": 0,
        "active_transport": "DISCONNECTED",
        "bridge_status": "OFFLINE",
        "q_value": 0.0
      };
    }
  }

  static Future<bool> setTransportOverride(String transport) async {
    try {
      await _channel.invokeMethod('setTransportOverride', {'transport': transport});
      return true;
    } on PlatformException {
      return false;
    }
  }

  static Future<String> exportLogs() async {
    try {
      final String path = await _channel.invokeMethod('exportLogs');
      return path;
    } on PlatformException {
      return "Failed to export logs";
    }
  }
}
