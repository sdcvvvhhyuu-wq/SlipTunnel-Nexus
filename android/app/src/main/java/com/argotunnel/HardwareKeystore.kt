package com.argotunnel

import android.content.Context
import android.content.SharedPreferences
import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import java.security.KeyStore
import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec

/**
 * HardwareKeystore implements Envelope Encryption.
 * It generates a hardware-backed AES key in the Android Secure Enclave,
 * which is used to encrypt/decrypt a volatile 32-byte master key passed to the Go EncryptedDB.
 */
class HardwareKeystore(private val context: Context) {

    private val keyAlias = "SlipTunnelMasterKeyAlias"
    private val prefs: SharedPreferences = context.getSharedPreferences("crypto_prefs", Context.MODE_PRIVATE)
    private val keyStore: KeyStore = KeyStore.getInstance("AndroidKeyStore").apply { load(null) }

    init {
        if (!keyStore.containsAlias(keyAlias)) {
            generateHardwareKey()
            val rawGoKey = ByteArray(32).apply { SecureRandom().nextBytes(this) }
            val encryptedGoKey = encryptWithHardwareKey(rawGoKey)
            prefs.edit().putString("encrypted_go_key", Base64.encodeToString(encryptedGoKey, Base64.DEFAULT)).apply()
        }
    }

    private fun generateHardwareKey() {
        val keyGenerator = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, "AndroidKeyStore")
        val keyGenParameterSpec = KeyGenParameterSpec.Builder(
            keyAlias,
            KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT
        )
            .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
            .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
            .setKeySize(256)
            .build()
        keyGenerator.init(keyGenParameterSpec)
        keyGenerator.generateKey()
    }

    private fun encryptWithHardwareKey(data: ByteArray): ByteArray {
        val secretKey = keyStore.getKey(keyAlias, null) as SecretKey
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        cipher.init(Cipher.ENCRYPT_MODE, secretKey)
        val iv = cipher.iv
        val ciphertext = cipher.doFinal(data)
        return iv + ciphertext
    }

    private fun decryptWithHardwareKey(encryptedData: ByteArray): ByteArray {
        val secretKey = keyStore.getKey(keyAlias, null) as SecretKey
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        val iv = encryptedData.copyOfRange(0, 12)
        val ciphertext = encryptedData.copyOfRange(12, encryptedData.size)
        val spec = GCMParameterSpec(128, iv)
        cipher.init(Cipher.DECRYPT_MODE, secretKey, spec)
        return cipher.doFinal(ciphertext)
    }

    fun getVolatileGoMasterKey(): String {
        val encryptedBase64 = prefs.getString("encrypted_go_key", null)
            ?: throw IllegalStateException("Master key not found")
        val encryptedData = Base64.decode(encryptedBase64, Base64.DEFAULT)
        val rawKey = decryptWithHardwareKey(encryptedData)
        
        // Convert to Hex string for Go C.CString consumption
        return rawKey.joinToString("") { "%02x".format(it) }
    }
}
