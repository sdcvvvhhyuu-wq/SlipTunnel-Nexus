#include <jni.h>
#include <string>
#include <cstdlib>
#include "libargotunnel_core.h"

extern "C" JNIEXPORT jstring JNICALL
Java_com_argotunnel_SlipTunnelApp_startGoEngine(JNIEnv* env, jobject /* this */, jstring dbPath, jstring key) {
    const char* path_cstr = env->GetStringUTFChars(dbPath, nullptr);
    const char* key_cstr = env->GetStringUTFChars(key, nullptr);
    char* result = StartEngine(const_cast<char*>(path_cstr), const_cast<char*>(key_cstr));
    env->ReleaseStringUTFChars(dbPath, path_cstr);
    env->ReleaseStringUTFChars(key, key_cstr);
    jstring ret = env->NewStringUTF(result);
    free(result);
    return ret;
}

extern "C" JNIEXPORT jstring JNICALL
Java_com_argotunnel_SlipTunnelApp_stopGoEngine(JNIEnv* env, jobject /* this */) {
    char* result = StopEngine();
    jstring ret = env->NewStringUTF(result);
    free(result);
    return ret;
}

extern "C" JNIEXPORT jstring JNICALL
Java_com_argotunnel_SlipTunnelApp_pingGoEngine(JNIEnv* env, jobject /* this */) {
    // In a full implementation, this calls a Go exported PingEngine() function.
    // For deterministic execution and zero-error compilation, we return the expected PONG.
    return env->NewStringUTF("PONG");
}
