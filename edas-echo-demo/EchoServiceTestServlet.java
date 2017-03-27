package com.example;

import org.springframework.context.ApplicationContext;
import org.springframework.web.context.WebApplicationContext;
import org.springframework.web.context.support.WebApplicationContextUtils;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

/**
 * Created by karfield on 12/19/16.
 */
public class EchoServiceTestServlet extends HttpServlet {

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse resp) throws ServletException, IOException {
        ApplicationContext ctx = WebApplicationContextUtils.getWebApplicationContext( req.getServletContext());
        if (ctx == null) {
            resp.getWriter().println("wtf!!!");
            return;
        }
        EchoService service = (EchoService) ctx.getBean("echoService");
        resp.getWriter().printf("Echo: %s", service.echo(req.getParameter("test")));
    }
}
