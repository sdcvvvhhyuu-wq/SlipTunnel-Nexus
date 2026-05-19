pub mod packet_runtime;

use jni::JNIEnv;
use jni::objects::JClass;
use jni::sys::jint;
use packet_runtime::TunOrchestrator;
use std::sync::Mutex;
use lazy_static::lazy_static;

lazy_static! {
    static ref TUN_ORCHESTRATOR: Mutex<Option<TunOrchestrator>> = Mutex::new(None);
}

// JNI Entry Point for Android VpnService (No-Root)
#[no_mangle]
pub extern "system" fn Java_com_argotunnel_SlipVpnService_startNativeRuntime(
    mut _env: JNIEnv,
    _class: JClass,
    fd: jint,
) -> jint {
    let mut orchestrator_guard = TUN_ORCHESTRATOR.lock().unwrap();
    
    let orchestrator = TunOrchestrator::new(fd);
    *orchestrator_guard = Some(orchestrator);
    
    0 // Success
}

#[no_mangle]
pub extern "system" fn Java_com_argotunnel_SlipVpnService_stopNativeRuntime(
    mut _env: JNIEnv,
    _class: JClass,
) -> jint {
    let mut orchestrator_guard = TUN_ORCHESTRATOR.lock().unwrap();
    if let Some(mut orchestrator) = orchestrator_guard.take() {
        orchestrator.shutdown();
    }
    0 // Success
}
