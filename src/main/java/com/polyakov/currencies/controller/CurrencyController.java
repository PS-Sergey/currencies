package com.polyakov.currencies.controller;

import com.polyakov.currencies.service.GiphyService;
import com.polyakov.currencies.service.OpenExchangeRatesService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RequestMapping("/api")
@RestController
public class CurrencyController {

    @Value("${giphy.rich}")
    private String rich;

    @Value("${giphy.broke}")
    private String broke;

    @Value("${giphy.strong}")
    private String strong;

    private final OpenExchangeRatesService openExchangeRatesService;
    private final GiphyService giphyService;

    @Autowired
    public CurrencyController(OpenExchangeRatesService openExchangeRatesService, GiphyService giphyService) {
        this.openExchangeRatesService = openExchangeRatesService;
        this.giphyService = giphyService;
    }

    @GetMapping("/currencies")
    public Map<String, String> getCurrencies() {
        return openExchangeRatesService.getCurrencies();
    }

    @GetMapping(value = "/rate", produces = MediaType.IMAGE_GIF_VALUE)
    public ResponseEntity<byte[]> getRate() {
        String tag = "";
        int counting = openExchangeRatesService.compareRate();
        switch (counting) {
            case 1:
                tag = rich;
                break;
            case -1:
                tag = broke;
                break;
            case 0:
                tag = strong;
                break;
        }
        byte[] gif = giphyService.getGif(tag);
        return ResponseEntity
                    .ok()
                    .contentType(MediaType.IMAGE_GIF)
                    .body(gif);
    }
}
