package com.argotunnel

import android.app.Application
import android.content.Intent
import android.util.Log
import java.io.File

class SlipTunnelApp : Application() {

    companion object {
        init {
            System.loadLibrary("argotunnel_core")      
            System.loadLibrary("argotunnel_rust_core") 
            System.loadLibrary("argotunnel_jni")       
        }
    }

    private lateinit var keystore: HardwareKeystore

    override fun onCreate() {
        super.onCreate()
        keystore = HardwareKeystore(this)
        
        val dbPath = File(filesDir, "nexus_secure.db").absolutePath
        val masterKey = keystore.getVolatileGoMasterKey()

        val result = startGoEngine(dbPath, masterKey)
        Log.i("SlipTunnelApp", "Go Engine Init: $result")

        // Initialize the Self-Healing Watchdog
        val watchdogIntent = Intent(this, WatchdogService::class.java)
        startService(watchdogIntent)
    }

    override fun onTerminate() {
        val result = stopGoEngine()
        Log.i("SlipTunnelApp", "Go Engine Shutdown: $result")
        
        val watchdogIntent = Intent(this, WatchdogService::class.java)
        stopService(watchdogIntent)
        
        super.onTerminate()
    }

    external fun startGoEngine(dbPath: String, key: String): String
    external fun stopGoEngine(): String
    external fun pingGoEngine(): String // Watchdog JNI Binding
}
