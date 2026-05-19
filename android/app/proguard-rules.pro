# =====================================================================
# PROJECT HIDDIFY-SLIPTUNNEL-NEXUS: MAXIMUM OBFUSCATION RULES
# =====================================================================

# 1. Preserve JNI FFI Boundaries (CRITICAL for Go/Rust/C++ interop)
-keepclasseswithmembernames class com.argotunnel.** {
    native <methods>;
}
-keep class com.argotunnel.SlipTunnelApp { *; }
-keep class com.argotunnel.SlipVpnService { *; }
-keep class com.argotunnel.WatchdogService { *; }

# 2. Preserve Flutter MethodChannels
-keep class io.flutter.** { *; }
-keep class io.flutter.plugin.common.** { *; }
-keep class com.argotunnel.MainActivity { *; }

# 3. Aggressive Obfuscation & Optimization
-repackageclasses 'com.argotunnel.core.obfuscated'
-allowaccessmodification
-optimizations !code/simplification/arithmetic,!field/*,!class/merging/*
-optimizationpasses 5

# 4. Strip Logging in Release Builds to prevent side-channel leakage
-assumenosideeffects class android.util.Log {
    public static boolean isLoggable(java.lang.String, int);
    public static int v(...);
    public static int i(...);
    public static int w(...);
    public static int d(...);
    public static int e(...);
}

# 5. Keep Android OS required attributes
-keepattributes *Annotation*
-keepattributes Signature
-keepattributes InnerClasses
