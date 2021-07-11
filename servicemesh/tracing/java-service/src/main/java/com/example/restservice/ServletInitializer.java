package com.example.restservice;

import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.boot.builder.SpringApplicationBuilder;
import org.springframework.boot.web.servlet.support.SpringBootServletInitializer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.context.request.RequestContextListener;
import org.springframework.web.filter.RequestContextFilter;

import javax.servlet.ServletContext;
import javax.servlet.ServletException;

@Configuration
@ComponentScan
@EnableAutoConfiguration
public class ServletInitializer extends SpringBootServletInitializer {

    @Override
    protected SpringApplicationBuilder configure(SpringApplicationBuilder application) {
        return application.sources(RestServiceApplication.class);
    }

//    @Bean
//    public RequestContextListener requestContextListener(){
//        return new RequestContextListener();
//    }

    @Override
    public void onStartup(ServletContext servletContext ) throws ServletException {
        super.onStartup(servletContext);
        servletContext.addFilter("requestContextFilter", new RequestContextFilter()).addMappingForUrlPatterns(null, false, "/*");
    }

}
