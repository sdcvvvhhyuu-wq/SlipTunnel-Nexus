package com.argotunnel

import android.content.Intent
import android.net.VpnService
import android.os.ParcelFileDescriptor
import android.util.Log
import java.io.IOException

class SlipVpnService : VpnService() {

    private var vpnInterface: ParcelFileDescriptor? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (vpnInterface != null) {
            Log.i("SlipVpnService", "VPN already running")
            return START_STICKY
        }

        try {
            val builder = Builder()
                .setSession("SlipTunnel Nexus")
                .addAddress("10.0.0.2", 24)
                .addDnsServer("1.1.1.1")
                .addRoute("0.0.0.0", 0)
                .setMtu(1500)
                .setBlocking(true)

            vpnInterface = builder.establish()

            vpnInterface?.let { pfd ->
                val fd = pfd.fd
                Log.i("SlipVpnService", "TUN Interface established. FD: $fd")
                
                // Pass the File Descriptor to the Rust Bare-Metal Runtime (Phase 2)
                val result = startNativeRuntime(fd)
                if (result != 0) {
                    Log.e("SlipVpnService", "Failed to start Rust Native Runtime")
                }
            }
        } catch (e: Exception) {
            Log.e("SlipVpnService", "Error establishing VPN", e)
        }

        return START_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        stopNativeRuntime()
        try {
            vpnInterface?.close()
            vpnInterface = null
        } catch (e: IOException) {
            Log.e("SlipVpnService", "Error closing VPN interface", e)
        }
    }

    // JNI Bindings directly to Rust (libargotunnel_rust_core.so)
    private external fun startNativeRuntime(fd: Int): Int
    private external fun stopNativeRuntime(): Int
}
