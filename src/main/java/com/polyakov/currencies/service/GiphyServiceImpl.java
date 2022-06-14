package com.polyakov.currencies.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.polyakov.currencies.exception.GifNotFoundException;
import com.polyakov.currencies.feign.GiphyFeign;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.RequestEntity;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class GiphyServiceImpl implements GiphyService{

    @Value("${giphy.api.key}")
    private String apiKey;

    private final GiphyFeign giphyFeign;
    private final RestTemplate restTemplate;
    private final ObjectMapper objectMapper;

    @Autowired
    public GiphyServiceImpl(GiphyFeign giphyFeign, RestTemplate restTemplate, ObjectMapper objectMapper) {
        this.giphyFeign = giphyFeign;
        this.restTemplate = restTemplate;
        this.objectMapper = objectMapper;
    }

    @Override
    public byte[] getGif(String tag) {
        String json = giphyFeign.getGif(apiKey, tag);
        String gifUrl = getGifUrl(json);
        RequestEntity requestEntity = RequestEntity
                                          .get(gifUrl)
                                          .build();
        ResponseEntity<byte[]> responseEntity = restTemplate.exchange(requestEntity, byte[].class);
        byte[] gif = responseEntity.getBody();
        if (gif == null) {
            throw new GifNotFoundException("Gif not found");
        }
        return gif;
    }

    private String getGifUrl(String json) {
        try {
            return objectMapper.readTree(json).get("data").get("images").get("original").get("url").asText();
        } catch (JsonProcessingException e) {
            throw new RuntimeException("Error during parsing response from giphy.com");
        }
    }
}
