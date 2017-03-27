package com.example;

import java.util.Random;

/**
 * Created by karfield on 12/19/16.
 */
public class EchoServiceProvider implements EchoService {
    public String echo(String str) {
        if (new Random().nextBoolean()) {
            throw new RuntimeException("WTF!");
        }
        return str;
    }
}
