package com.polyakov.currencies.service;

import com.polyakov.currencies.exception.GifNotFoundException;
import com.polyakov.currencies.feign.GiphyFeign;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.core.io.ClassPathResource;
import org.springframework.http.RequestEntity;
import org.springframework.http.ResponseEntity;
import org.springframework.web.client.RestTemplate;

import java.io.*;
import java.util.stream.Collectors;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.anyString;

@SpringBootTest
class GiphyServiceImplTest {

    @MockBean
    private GiphyFeign giphyFeign;

    @MockBean
    private RestTemplate restTemplate;

    @Autowired
    private GiphyServiceImpl giphyService;

    @Test
    void whenGetGifThenReturnByteArray() throws IOException {
        InputStream resource = new ClassPathResource("parseTest").getInputStream();
        BufferedReader reader = new BufferedReader(new InputStreamReader(resource));
        String giphyFeignResponse = reader
                                        .lines()
                                        .collect(Collectors.joining("\n"));
        String gifUrl = "TestUrl";
        RequestEntity requestEntity = RequestEntity
                                          .get(gifUrl)
                                          .build();
        byte[] gif = new byte[1];

        Mockito.when(giphyFeign.getGif(anyString(), anyString()))
            .thenReturn(giphyFeignResponse);
        Mockito.when(restTemplate.exchange(requestEntity, byte[].class))
            .thenReturn(ResponseEntity.ok().body(gif));
        byte[] result = giphyService.getGif("testTag");
        assertEquals(gif, result);
    }

    @Test
    void whenGifUrlReturnNullThrowException() throws IOException {
        InputStream resource = new ClassPathResource("parseTest").getInputStream();
        BufferedReader reader = new BufferedReader(new InputStreamReader(resource));
        String giphyFeignResponse = reader
                .lines()
                .collect(Collectors.joining("\n"));
        String gifUrl = "TestUrl";
        RequestEntity requestEntity = RequestEntity
                .get(gifUrl)
                .build();

        Mockito.when(giphyFeign.getGif(anyString(), anyString()))
                .thenReturn(giphyFeignResponse);
        Mockito.when(restTemplate.exchange(requestEntity, byte[].class))
                .thenReturn(ResponseEntity.ok().body(null));
        Throwable thrown = assertThrows(GifNotFoundException.class, () -> {
            giphyService.getGif("testTag");
        });
        assertNotNull(thrown.getMessage());
    }
}