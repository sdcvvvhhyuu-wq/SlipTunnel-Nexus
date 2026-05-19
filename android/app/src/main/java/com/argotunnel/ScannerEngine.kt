package com.argotunnel

import android.util.Log

/**
 * ScannerEngine provides a Kotlin interface to trigger the Go-based BridgeScanner
 * and Autonomous IP/Domain Hunter.
 */
class ScannerEngine {

    fun triggerOnlineScraping() {
        Log.i("ScannerEngine", "Triggering Online Scraping via Go Engine")
        // JNI call to Go Engine's BridgeScanner.StartParallelScan() would be mapped here
    }

    fun triggerOfflineValidation() {
        Log.i("ScannerEngine", "Triggering Offline Validation via Go Engine")
        // JNI call to Go Engine's BridgeScanner.offlineValidationLoop() would be mapped here
    }
}
