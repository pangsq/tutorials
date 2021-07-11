package com.example.restservice;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class ForwardService {
    @Autowired
    private TracingContext context;

    public String forward(String targetURL) {
        //        HttpHeaders httpHeaders = Tracing.injectHeaders(new HttpHeaders()); // it uses the thread-scoped context
        HttpHeaders httpHeaders = context.inject(new HttpHeaders()); // it uses the request-scoped context
//        System.out.println(httpHeaders);

        return new RestTemplate().exchange(targetURL, HttpMethod.GET,
                new HttpEntity("parameters", httpHeaders), String.class).getBody();
    }

    public String forwardWithoutTracing(String targetURL) {
        HttpHeaders httpHeaders = new HttpHeaders();

        return new RestTemplate().exchange(targetURL, HttpMethod.GET,
                new HttpEntity("parameters", httpHeaders), String.class).getBody();
    }
}
