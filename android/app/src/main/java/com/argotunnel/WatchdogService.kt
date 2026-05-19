package com.argotunnel

import android.app.Service
import android.content.Intent
import android.os.IBinder
import android.util.Log
import kotlinx.coroutines.*

/**
 * WatchdogService runs in a detached coroutine scope.
 * It continuously pings the Go/Rust runtime via JNI. If the engine hangs, deadlocks,
 * or crashes, it autonomously tears down the VPN interface and hot-reloads the subsystems.
 */
class WatchdogService : Service() {

    private val job = SupervisorJob()
    private val scope = CoroutineScope(Dispatchers.IO + job)
    private var isRunning = false

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        if (!isRunning) {
            isRunning = true
            Log.i("WatchdogService", "Self-Healing Watchdog Initialized")
            startWatchdogLoop()
        }
        return START_STICKY
    }

    private fun startWatchdogLoop() {
        scope.launch {
            while (isActive) {
                delay(10000) // Health check interval: 10 seconds
                try {
                    val app = application as SlipTunnelApp
                    val response = app.pingGoEngine()
                    
                    if (response != "PONG") {
                        Log.e("WatchdogService", "Engine returned invalid state: $response")
                        triggerRecovery()
                    }
                } catch (e: Exception) {
                    Log.e("WatchdogService", "Engine unresponsive (Deadlock/Crash detected)", e)
                    triggerRecovery()
                }
            }
        }
    }

    private fun triggerRecovery() {
        Log.w("WatchdogService", "EXECUTING EMERGENCY SELF-HEALING RECOVERY...")
        
        // 1. Tear down the hung VPN Service
        val stopIntent = Intent(this@WatchdogService, SlipVpnService::class.java)
        stopService(stopIntent)

        // 2. Allow OS to reclaim file descriptors and memory
        Thread.sleep(2500) 

        // 3. Hot-reload the VPN Service
        val startIntent = Intent(this@WatchdogService, SlipVpnService::class.java)
        startService(startIntent)
        
        Log.i("WatchdogService", "Recovery Complete. Subsystems Online.")
    }

    override fun onDestroy() {
        super.onDestroy()
        job.cancel()
        isRunning = false
    }
}
