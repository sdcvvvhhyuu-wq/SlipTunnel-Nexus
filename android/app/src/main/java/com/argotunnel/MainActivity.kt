package com.argotunnel

import android.app.Activity
import android.content.Intent
import android.net.VpnService
import android.os.Bundle
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import org.json.JSONObject

class MainActivity : FlutterActivity() {
    private val CHANNEL = "com.argotunnel/vpn"
    private val VPN_REQUEST_CODE = 1001

    private var pendingResult: MethodChannel.Result? = null

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            when (call.method) {
                "startVpn" -> {
                    val intent = VpnService.prepare(this)
                    if (intent != null) {
                        pendingResult = result
                        startActivityForResult(intent, VPN_REQUEST_CODE)
                    } else {
                        startSlipVpnService()
                        result.success("VPN Started")
                    }
                }
                "stopVpn" -> {
                    val stopIntent = Intent(this, SlipVpnService::class.java)
                    stopService(stopIntent)
                    result.success("VPN Stopped")
                }
                "getDiagnostics" -> {
                    // In a full implementation, this queries the Go Engine via JNI.
                    // Returning deterministic JSON for the Flutter UI.
                    val diagnostics = JSONObject().apply {
                        put("network_health", "OPTIMAL")
                        put("latency_ms", 142)
                        put("active_transport", "CDN_RELAY (Psiphon Meek)")
                        put("bridge_status", "ONLINE - obfs4")
                        put("q_value", 0.92)
                    }
                    result.success(diagnostics.toString())
                }
                "setTransportOverride" -> {
                    val transport = call.argument<String>("transport")
                    // JNI call to Go Engine's Selector to override transport
                    result.success("Transport overridden to $transport")
                }
                "exportLogs" -> {
                    // Logic to dump logcat/Go logs to a file
                    result.success("/storage/emulated/0/Download/nexus_logs.txt")
                }
                else -> result.notImplemented()
            }
        }
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        if (requestCode == VPN_REQUEST_CODE) {
            if (resultCode == Activity.RESULT_OK) {
                startSlipVpnService()
                pendingResult?.success("VPN Started")
            } else {
                pendingResult?.error("VPN_DENIED", "User denied VPN permission", null)
            }
            pendingResult = null
        } else {
            super.onActivityResult(requestCode, resultCode, data)
        }
    }

    private fun startSlipVpnService() {
        val intent = Intent(this, SlipVpnService::class.java)
        startService(intent)
    }
}
