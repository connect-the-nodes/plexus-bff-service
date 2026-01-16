package com.zynchub.digital.hubservice.app.feature.flags.service;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import com.zynchub.digital.hubservice.app.feature.flags.exception.FeatureNotFoundException;
import com.zynchub.digital.hubservice.app.feature.flags.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.feature.flags.model.FeatureFlag;
import com.zynchub.digital.hubservice.app.feature.flags.service.impl.PropertyFeaturesRetrieverImpl;
import java.util.List;
import org.junit.jupiter.api.Test;

class FeaturesRetrieverTest {

    @Test
    void isActive_returns_true_when_parent_chain_enabled() {
        FeaturesRetriever retriever =
                () ->
                        List.of(
                                buildFlag("PARENT", "PARENT", true),
                                buildFlag("CHILD", "PARENT", true));

        assertThat(retriever.isActive("CHILD")).isTrue();
    }

    @Test
    void isActive_returns_false_when_parent_disabled() {
        FeaturesRetriever retriever =
                () ->
                        List.of(
                                buildFlag("PARENT", "PARENT", false),
                                buildFlag("CHILD", "PARENT", true));

        assertThat(retriever.isActive("CHILD")).isFalse();
    }

    @Test
    void isActive_throws_when_feature_missing() {
        FeaturesRetriever retriever = List::of;

        assertThatThrownBy(() -> retriever.isActive("MISSING"))
                .isInstanceOf(FeatureNotFoundException.class);
    }

    @Test
    void propertyRetriever_reads_yaml_from_classpath() {
        FeaturesRetriever retriever =
                new PropertyFeaturesRetrieverImpl("test-features.yml", new FeatureMapper());

        assertThat(retriever.isActive("TEST_PARENT")).isTrue();
        assertThat(retriever.isActive("TEST_CHILD")).isTrue();
    }

    private FeatureFlag buildFlag(String name, String parent, boolean enabled) {
        FeatureFlag flag = new FeatureFlag();
        flag.setName(name);
        flag.setParent(parent);
        flag.setEnabled(enabled);
        return flag;
    }
}
