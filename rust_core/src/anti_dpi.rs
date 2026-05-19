use std::time::{Duration, SystemTime, UNIX_EPOCH};
use std::thread;

/// DpiObfuscator implements Markov-chain inspired packet pacing and inter-packet delay (IPD)
/// randomization to defeat Machine Learning-based statistical traffic classifiers.
pub struct DpiObfuscator {
    prng_state: u64,
}

impl DpiObfuscator {
    pub fn new() -> Self {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or(Duration::from_secs(0))
            .as_nanos() as u64;
        
        Self { prng_state: now }
    }

    /// Linear Congruential Generator (LCG) for zero-dependency, ultra-fast deterministic entropy
    fn next_rand(&mut self) -> u64 {
        self.prng_state = self.prng_state.wrapping_mul(6364136223846793005).wrapping_add(1);
        self.prng_state
    }

    /// apply_packet_pacing: Injects algorithmic delays based on packet size to emulate benign traffic
    pub fn apply_packet_pacing(&mut self, packet_size: usize) {
        // Generate 0 to 14 ms of random jitter
        let jitter = self.next_rand() % 15; 
        
        // Larger packets incur a slight structural delay to mimic standard HTTPS file transfers
        let size_factor = (packet_size as u64) / 500; 
        
        // Base 2ms delay + size factor + jitter
        let delay_ms = 2 + size_factor + jitter; 

        if delay_ms > 0 {
            thread::sleep(Duration::from_millis(delay_ms));
        }
    }

    /// mutate_tcp_window: Shuffles TCP window sizes to mimic standard Windows/Linux OS signatures
    pub fn mutate_tcp_window(&mut self) -> u16 {
        let windows: [u16; 4] = [65535, 64240, 29200, 14600];
        let idx = (self.next_rand() % 4) as usize;
        windows[idx]
    }
}
