package com.example.restservice;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpHeaders;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
public class ForwardController {

    @Autowired
    private TracingContext context;

    @Autowired
    private ForwardService service;

    @Value("${NEXT_TARGET:http://go-service}")
    private String nextTarget;

    @Value("${FINAL_TARGET:http://httpbin}")
    private String finalTarget;

    @GetMapping("{path}")
    public Object forward(@RequestHeader HttpHeaders headers, @PathVariable(value = "path") String path, @RequestParam(value = "ttl", defaultValue = "0") String ttl) {
//        Tracing.extract(headers); // it will extract headers and store the context in current thread (scope = thread)
//        context.extract(headers); // it will extract headers and store the context in current request (scope = request)
        Integer ttlNum = 0;
        try {
            ttlNum = Integer.valueOf(ttl);
        } catch (NumberFormatException nfe) {
            System.err.println(nfe);
        }
        // for debug
//        {
//            System.out.println(headers);
//            List<String> traceIds = headers.get("X-B3-TraceId");
//            List<String> spanIds = headers.get("X-B3-SpanId");
//            if (traceIds != null && !traceIds.isEmpty() && spanIds != null && !spanIds.isEmpty()) {
//                System.out.printf("b3: %s/%s\n", traceIds.get(0), spanIds.get(0));
//            }
//        }
//        System.out.printf("ttl=%d\n", ttlNum);
        String targetURL = "";
        if (ttlNum == 0) {
            targetURL = finalTarget + "/" + path;
        } else {
            targetURL = nextTarget + "/" + path + "?ttl=" + (ttlNum - 1);
        }
        return service.forward(targetURL);
    }
}
