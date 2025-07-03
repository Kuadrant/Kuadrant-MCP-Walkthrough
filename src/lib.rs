use proxy_wasm::traits::*;
use proxy_wasm::types::*;

#[derive(Debug)]
struct DebugRoot;

impl Context for DebugRoot {}

impl RootContext for DebugRoot {
    fn on_vm_start(&mut self, _vm_configuration_size: usize) -> bool {
        proxy_wasm::hostcalls::log(LogLevel::Info, "[WASM] Debug WASM filter started").unwrap();
        true
    }

    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }

    fn create_http_context(&self, _context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(DebugHttpFilter))
    }
}

struct DebugHttpFilter;

impl Context for DebugHttpFilter {}

impl HttpContext for DebugHttpFilter {
    fn on_http_request_headers(&mut self, num_headers: usize, end_of_stream: bool) -> Action {
        proxy_wasm::hostcalls::log(
            LogLevel::Info,
            &format!("[WASM] Starting request header processing... (headers: {}, end_of_stream: {})", num_headers, end_of_stream)
        ).unwrap();
        
        proxy_wasm::hostcalls::log(LogLevel::Info, "[WASM] Completed request header processing").unwrap();
        Action::Continue
    }

    fn on_http_request_body(&mut self, body_size: usize, end_of_stream: bool) -> Action {
        proxy_wasm::hostcalls::log(
            LogLevel::Info,
            &format!("[WASM] Starting request body processing... (size: {}, end_of_stream: {})", body_size, end_of_stream)
        ).unwrap();
        
        if end_of_stream {
            proxy_wasm::hostcalls::log(LogLevel::Info, "[WASM] Completed request body processing").unwrap();
        }
        
        Action::Continue
    }

    fn on_http_response_headers(&mut self, num_headers: usize, end_of_stream: bool) -> Action {
        proxy_wasm::hostcalls::log(
            LogLevel::Info,
            &format!("[WASM] Starting response header processing... (headers: {}, end_of_stream: {})", num_headers, end_of_stream)
        ).unwrap();
        
        proxy_wasm::hostcalls::log(LogLevel::Info, "[WASM] Completed response header processing").unwrap();
        Action::Continue
    }

    fn on_http_response_body(&mut self, body_size: usize, end_of_stream: bool) -> Action {
        proxy_wasm::hostcalls::log(
            LogLevel::Info,
            &format!("[WASM] Starting response body processing... (size: {}, end_of_stream: {})", body_size, end_of_stream)
        ).unwrap();
        
        if end_of_stream {
            proxy_wasm::hostcalls::log(LogLevel::Info, "[WASM] Completed response body processing").unwrap();
        }
        
        Action::Continue
    }
}

proxy_wasm::main! {{
    proxy_wasm::set_log_level(LogLevel::Info);
    proxy_wasm::set_root_context(|_| -> Box<dyn RootContext> {
        Box::new(DebugRoot)
    });
}} 