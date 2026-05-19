package com.argotunnel

import android.app.Application
import android.util.Log
import java.io.File

class SlipTunnelApp : Application() {

    companion object {
        init {
            System.loadLibrary("argotunnel_core")      // Go Core
            System.loadLibrary("argotunnel_rust_core") // Rust Core
            System.loadLibrary("argotunnel_jni")       // C++ Bridge
        }
    }

    private lateinit var keystore: HardwareKeystore

    override fun onCreate() {
        super.onCreate()
        keystore = HardwareKeystore(this)
        
        val dbPath = File(filesDir, "nexus_secure.db").absolutePath
        val masterKey = keystore.getVolatileGoMasterKey()

        // Initialize the Go Orchestrator autonomously
        val result = startGoEngine(dbPath, masterKey)
        Log.i("SlipTunnelApp", "Go Engine Init: $result")
    }

    override fun onTerminate() {
        val result = stopGoEngine()
        Log.i("SlipTunnelApp", "Go Engine Shutdown: $result")
        super.onTerminate()
    }

    // JNI Bindings to C++ Wrapper (which calls Go)
    external fun startGoEngine(dbPath: String, key: String): String
    external fun stopGoEngine(): String
}
