plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "com.argotunnel"
    compileSdk = 34

    defaultConfig {
        applicationId = "com.argotunnel"
        minSdk = 24
        targetSdk = 34
        versionCode = 1
        versionName = "1.4.0-OMNI-FINAL"

        externalNativeBuild {
            cmake {
                cppFlags += "-std=c++17"
            }
        }
        ndk {
            abiFilters.add("arm64-v8a")
        }
    }

    // ─────────────────────────────────────────────────────────────────────────
    // SIGNING CONFIG — Kotlin DSL
    // اگر متغیرهای محیطی تنظیم شده باشند (CI)، از آن‌ها استفاده می‌شود.
    // در غیر این صورت، debug keystore پیش‌فرض Android SDK به‌کار می‌رود.
    // ─────────────────────────────────────────────────────────────────────────
    signingConfigs {
        create("release") {
            val envStoreFile = System.getenv("SIGNING_STORE_FILE")
            if (!envStoreFile.isNullOrEmpty()) {
                storeFile     = file(envStoreFile)
                storePassword = System.getenv("SIGNING_STORE_PASSWORD")
                keyAlias      = System.getenv("SIGNING_KEY_ALIAS")
                keyPassword   = System.getenv("SIGNING_KEY_PASSWORD")
            } else {
                // Fallback: debug keystore همیشه روی SDK موجود است
                storeFile     = file("${System.getProperty("user.home")}/.android/debug.keystore")
                storePassword = "android"
                keyAlias      = "androiddebugkey"
                keyPassword   = "android"
            }
        }
    }

    buildTypes {
        release {
            signingConfig    = signingConfigs.getByName("release")
            isMinifyEnabled  = true
            isShrinkResources = true
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }

    externalNativeBuild {
        cmake {
            path    = file("src/main/cpp/CMakeLists.txt")
            version = "3.22.1"
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }
    kotlinOptions {
        jvmTarget = "17"
    }
}

dependencies {
    implementation("androidx.core:core-ktx:1.12.0")
    implementation("androidx.appcompat:appcompat:1.6.1")
    implementation("com.google.android.material:material:1.11.0")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-android:1.7.3")
}
